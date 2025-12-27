import { defineConfig } from 'vocs'

export default defineConfig({
  basePath: '/rollups-golang-course',
  baseUrl: 'https://henriquemarlon.github.io',
  description: 'Cartesi Rollups Golang Course',
  title: 'Docs',
  sidebar: [
    {
      text: 'Getting Started',
      link: '/getting-started',
    },
    {
      text: 'Example',
      link: '/example',
    },
  ],
})