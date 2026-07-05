package core

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/nyaruka/phonenumbers"
	"sptzx/src/config"
)

var strictPhonePattern = regexp.MustCompile(`^[1-9][0-9]{6,14}$`)

func ResolveLoginConfigInteractive(cfg *config.Config) error {
	if cfg == nil || !isInteractiveTerminal() {
		return nil
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprintln(os.Stdout, "WhatsApp Login Setup")
	fmt.Fprintln(os.Stdout, "Pilih metode login:")
	fmt.Fprintln(os.Stdout, "1. QR")
	fmt.Fprintln(os.Stdout, "2. Pairing Code")

	for {
		fmt.Fprint(os.Stdout, "Masukkan pilihan (1/2): ")
		choice, err := readLine(reader)
		if err != nil {
			return fmt.Errorf("gagal membaca pilihan login: %w", err)
		}

		switch normalizeChoice(choice) {
		case "qr":
			cfg.LoginMethod = "qr"
			cfg.PairingPhone = ""
			return nil
		case "paircode":
			cfg.LoginMethod = "paircode"

			for {
				fmt.Fprint(os.Stdout, "Masukkan nomor internasional dengan kode negara (contoh: 6281234567890): ")
				raw, err := readLine(reader)
				if err != nil {
					return fmt.Errorf("gagal membaca nomor telepon: %w", err)
				}

				normalized, err := normalizePairingPhone(raw)
				if err != nil {
					fmt.Fprintf(os.Stdout, "Input tidak valid: %s\n", err.Error())
					continue
				}

				cfg.PairingPhone = normalized
				return nil
			}
		default:
			fmt.Fprintln(os.Stdout, "Pilihan tidak valid. Gunakan 1 untuk QR atau 2 untuk Pairing Code.")
		}
	}
}

func normalizeChoice(input string) string {
	v := strings.ToLower(strings.TrimSpace(input))
	switch v {
	case "1", "qr":
		return "qr"
	case "2", "pair", "pairing", "paircode":
		return "paircode"
	default:
		return ""
	}
}

func normalizePairingPhone(input string) (string, error) {
	v := strings.TrimSpace(input)
	if !strictPhonePattern.MatchString(v) {
		return "", fmt.Errorf("nomor harus angka saja dan harus mulai dengan kode negara")
	}

	num, err := phonenumbers.Parse("+"+v, "ZZ")
	if err != nil {
		return "", fmt.Errorf("format nomor tidak valid")
	}

	if !phonenumbers.IsValidNumber(num) {
		return "", fmt.Errorf("nomor tidak valid menurut aturan nomor internasional")
	}

	e164 := phonenumbers.Format(num, phonenumbers.E164)
	return strings.TrimPrefix(e164, "+"), nil
}

func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		if errors.Is(err, io.EOF) {
			return strings.TrimSpace(line), nil
		}
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func isInteractiveTerminal() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
