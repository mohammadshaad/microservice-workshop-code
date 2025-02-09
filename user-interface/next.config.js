/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'standalone',
  poweredByHeader: false,
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: 'http://localhost/api/:path*',
      },
    ];
  },
  // Add assetPrefix for production
  assetPrefix: process.env.NODE_ENV === 'production' ? 'http://localhost' : '',
}

module.exports = nextConfig 