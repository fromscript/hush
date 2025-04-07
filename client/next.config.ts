// next.config.js
import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: 'standalone',
  publicRuntimeConfig: {
    wsUrl: process.env.NEXT_PUBLIC_WS_URL
  },
  async headers() {
    return [
      {
        source: '/_next/static/(.*)',
        headers: [
          {
            key: 'Access-Control-Allow-Origin',
            value: process.env.CORS_ALLOWED_ORIGINS || '*',
          },
        ],
      },
    ]
  },
};

export default nextConfig;
