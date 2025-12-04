// next.config.ts
import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  images: {
    domains: ['localhost'],
    unoptimized: process.env.NODE_ENV === 'development',
  },
  // Уберите experimental если он был
  // Добавьте компилятор для оптимизации
  compiler: {
    removeConsole: process.env.NODE_ENV === 'production',
  },
  // Уберите кастомные headers если они не критичны
  async headers() {
    return [
      {
        source: '/:path*',
        headers: [
          { 
            key: 'Access-Control-Allow-Origin', 
            value: process.env.NODE_ENV === 'development' 
              ? '*' 
              : 'https://your-production-domain.com' 
          },
        ],
      },
    ]
  },
  // Добавьте для отладки
  logging: {
    fetches: {
      fullUrl: true,
    },
  },
}

export default nextConfig;