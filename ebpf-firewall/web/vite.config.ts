import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'
import tsconfigPaths from 'vite-tsconfig-paths'
import { vitePluginForArco } from '@arco-plugins/vite-react'

export default defineConfig({
	plugins: [react(), tsconfigPaths(), vitePluginForArco({ style: 'css' })],
	server: {
		port: 8080
	}
})
