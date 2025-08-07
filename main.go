package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/crypto/curve25519"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Kullanım: wges create [server|client] [name]")
		os.Exit(1)
	}

	action, typ, name := os.Args[1], os.Args[2], os.Args[3]
	if action != "create" {
		fmt.Println("Geçersiz komut. Örnek: wges create server wg0")
		os.Exit(1)
	}

	switch typ {
	case "server":
		createServer(name)
	case "client":
		createClient(name)
	default:
		fmt.Println("Geçersiz tip. server veya client olmalı.")
		os.Exit(1)
	}
}

// Curve25519 tabanlı WireGuard anahtar çifti üretir
func generateKeyPair() (string, string, error) {
	var privKey [32]byte
	if _, err := rand.Read(privKey[:]); err != nil {
		return "", "", err
	}
	var pubKey [32]byte
	curve25519.ScalarBaseMult(&pubKey, &privKey)

	privB64 := base64.StdEncoding.EncodeToString(privKey[:])
	pubB64 := base64.StdEncoding.EncodeToString(pubKey[:])

	return privB64, pubB64, nil
}

// IP adresini bir artırır (örnek: 10.8.0.2 -> 10.8.0.3)
func incrementIP(ip string) string {
	var a, b, c, d int
	fmt.Sscanf(ip, "%d.%d.%d.%d", &a, &b, &c, &d)
	d++
	if d > 254 {
		d = 2
		c++
	}
	return fmt.Sprintf("%d.%d.%d.%d", a, b, c, d)
}

func createServer(wgName string) {
	fmt.Printf("Sunucu oluşturuluyor: %s\n", wgName)

	priv, pub, err := generateKeyPair()
	if err != nil {
		fmt.Println("Anahtar üretim hatası:", err)
		os.Exit(1)
	}

	wgDir := "/etc/wireguard"
	if err := os.MkdirAll(wgDir, 0700); err != nil {
		fmt.Println("Dizin oluşturulamadı:", err)
		os.Exit(1)
	}

	// Anahtarları dosyalara yaz
	_ = os.WriteFile(filepath.Join(wgDir, wgName+".key"), []byte(priv), 0600)
	_ = os.WriteFile(filepath.Join(wgDir, wgName+".pub"), []byte(pub), 0644)

	// IP durumu ve server bilgileri
	serverIP := "10.8.0.1/24"
	_ = os.WriteFile(filepath.Join(wgDir, "ip_state"), []byte("10.8.0.2"), 0600)
	info := fmt.Sprintf("%s\n%s\n", pub, "10.8.0.1")
	_ = os.WriteFile(filepath.Join(wgDir, "server_info"), []byte(info), 0644)

	// Sunucu yapılandırması
	conf := fmt.Sprintf(`[Interface]
PrivateKey = %s
Address = %s
ListenPort = 51820
`, priv, serverIP)
	_ = os.WriteFile(filepath.Join(wgDir, wgName+".conf"), []byte(conf), 0600)

	fmt.Println("Sunucu yapılandırması oluşturuldu:")
	fmt.Println("  Private Key:", filepath.Join(wgDir, wgName+".key"))
	fmt.Println("  Public Key :", filepath.Join(wgDir, wgName+".pub"))
	fmt.Println("  Config     :", filepath.Join(wgDir, wgName+".conf"))
}

func createClient(clientName string) {
	fmt.Printf("Client oluşturuluyor: %s\n", clientName)

	priv, pub, err := generateKeyPair()
	if err != nil {
		fmt.Println("Anahtar üretim hatası:", err)
		os.Exit(1)
	}

	clientDir := "/etc/wireguard/clients"
	if err := os.MkdirAll(clientDir, 0700); err != nil {
		fmt.Println("Client dizini oluşturulamadı:", err)
		os.Exit(1)
	}

	_ = os.WriteFile(filepath.Join(clientDir, clientName+".key"), []byte(priv), 0600)
	_ = os.WriteFile(filepath.Join(clientDir, clientName+".pub"), []byte(pub), 0644)

	wgDir := "/etc/wireguard"
	ipStatePath := filepath.Join(wgDir, "ip_state")
	lastIPBytes, err := os.ReadFile(ipStatePath)
	if err != nil {
		fmt.Println("IP state okunamadı:", err)
		os.Exit(1)
	}
	lastIP := string(lastIPBytes)
	nextIP := incrementIP(lastIP)
	_ = os.WriteFile(ipStatePath, []byte(nextIP), 0600)
	clientIP := lastIP + "/32"

	infoBytes, err := os.ReadFile(filepath.Join(wgDir, "server_info"))
	if err != nil {
		fmt.Println("Server info okunamadı:", err)
		os.Exit(1)
	}
	var serverPub, serverIP string
	fmt.Sscanf(string(infoBytes), "%s\n%s\n", &serverPub, &serverIP)

	conf := fmt.Sprintf(`[Interface]
PrivateKey = %s
Address = %s

[Peer]
PublicKey = %s
Endpoint = %s:51820
AllowedIPs = 0.0.0.0/0
`, priv, clientIP, serverPub, serverIP)
	_ = os.WriteFile(filepath.Join(clientDir, clientName+".conf"), []byte(conf), 0600)

	fmt.Println("Client yapılandırması oluşturuldu:")
	fmt.Println("  Private Key:", filepath.Join(clientDir, clientName+".key"))
	fmt.Println("  Public Key :", filepath.Join(clientDir, clientName+".pub"))
	fmt.Println("  Config     :", filepath.Join(clientDir, clientName+".conf"))
	fmt.Println("  Atanan IP  :", clientIP)
}
