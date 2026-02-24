/** @type {import('next').NextConfig} */
const nextConfig = {
    devIndicators: {
        autoPrerender: false
    },
    webpackDevMiddleware: config => {
        config.watchOptions = {
            poll: 1000,
            aggregateTimeout: 300
        }
        return config
    },
    experimental: {
        hmrHost: '192.168.211.128'
    }
};

module.exports = nextConfig;
