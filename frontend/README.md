# Toko Berkah POS - Frontend

Frontend untuk sistem kasir Toko Berkah menggunakan Vue.js 3 + TypeScript + Tailwind CSS.

## Setup

### Prerequisites
- Node.js 18+
- npm atau yarn

### Installation

```bash
cd frontend
npm install
```

### Environment Variables

Copy `.env.example` ke `.env`:

```bash
cp .env.example .env
```

### Development

```bash
npm run dev
```

Frontend akan berjalan di `http://localhost:3000`

### Build

```bash
npm run build
```

Output akan berada di folder `dist/`

## Project Structure

- `src/main.ts` - Entry point
- `src/App.vue` - Root component
- `src/api/` - API calls
- `src/stores/` - Pinia state management
- `src/views/` - Page components
- `src/components/` - Reusable components
- `src/router/` - Route definitions
- `src/types/` - TypeScript types
- `src/utils/` - Utility functions
- `src/composables/` - Composable hooks
- `src/assets/` - Static assets

## Features

### Kasir View
- Input produk via barcode/SKU
- Management keranjang belanja
- Harga dinamis (retail vs grosir)
- Proses pembayaran (tunai/transfer)
- Cetak struk

### Owner Dashboard
- Overview penjualan harian
- Breakdown retail vs grosir
- Daftar produk terlaris
- Manajemen inventori
- Management kasir
- Export laporan
