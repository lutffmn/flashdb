# 🚀 FlashDB: Concurrent In-Memory Key-Value Store

**FlashDB** adalah library penyimpanan data _key-value_ sederhana berbasis memori (RAM) yang dibangun menggunakan bahasa Go. Project ini dirancang untuk menangani akses data intensif dari ribuan Goroutine secara bersamaan dengan efisiensi tinggi dan keamanan memori yang terjamin.

## 🛠️ Fitur Utama

- **Thread-Safe Operations**: Implementasi penuh `sync.RWMutex` untuk mencegah _race condition_.
- **High Performance Reads**: Optimasi menggunakan `RLock` yang memungkinkan akses baca paralel oleh banyak Goroutine secara bersamaan.
- **Fine-Grained Locking**: Pemisahan kunci (_lock_) antara data utama dan statistik metrik untuk meminimalkan _lock contention_.
- **Modern Go Standard**: Menggunakan pola konkurensi terbaru Go (termasuk fitur `wg.Go` dari v1.25.0+).

---

## 🏗️ Arsitektur Konkurensi

Project ini menerapkan strategi **Dual-Locking**:

1. **Data Safety (`sync.RWMutex`)**:
   - Melindungi integritas map `store`.
   - `RLock()` (Read Lock) digunakan pada operasi `Get` untuk performa maksimal.
   - `Lock()` (Write Lock) digunakan pada `Set` dan `Delete` untuk modifikasi data yang aman.
2. **Metrics Safety (`sync.Mutex`)**:
   - Digunakan secara terpisah untuk mencatat statistik (`totalReads`, `totalWrites`).
   - Hal ini memastikan operasi tulis statistik tidak memperlambat operasi baca data utama.

---

## 📦 Cara Penggunaan

```go
package main

import (
    "fmt"
    "[github.com/lutffmn/flashdb](https://github.com/lutffmn/flashdb)"
)

func main() {
    db := &FlashDB{
        store: make(map[string]DataItem),
    }

    db.Set("user_101", 1250)

    val, ok := db.Get("user_101")
    if ok {
        fmt.Println("Data ditemukan:", val)
    }

    reads, writes := db.GetStats()
    fmt.Printf("Total Reads: %d, Total Writes: %d\n", reads, writes)

```
