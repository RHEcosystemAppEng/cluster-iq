# ClusterIQ Console

![Version](https://img.shields.io/badge/version-0.4.2-blue) ![License](https://img.shields.io/badge/license-Apache--2.0-green)

---

ClusterIQ Console provides a Web UI for the [ClusterIQ](https://github.com/RHEcosystemAppEng/cluster-iq) project.

## Deployment

This section explains how to deploy ClusterIQ Console.

### Prerequisites

- [Node.js](https://nodejs.org/) 18.x or higher
- [npm](https://www.npmjs.com/) 8.x or higher

### Quick-start (local)

```sh
git clone git@github.com:RHEcosystemAppEng/cluster-iq-console.git
cd cluster-iq-console
npm install && npm run start
```

## Development scripts

```sh
# Install development/build dependencies
npm install

# Start the development server
npm run start

# Run a production build (outputs to "dist" dir)
npm run build

# Run the preview
npm run preview
```

## Configurations

- [TypeScript Config](./tsconfig.json)
- [Editor Config](./.editorconfig)
- [Prettier Config](./.prettierrc.js)

## Code quality tools

- To keep our code formatting in check, we use [prettier](https://github.com/prettier/prettier)
- Code formatting and validation is done with [Husky](https://typicode.github.io/husky) and [lint-staged](https://github.com/lint-staged/lint-staged)

## Getting Started with Vite

This project was bootstrapped with [Vite](https://vite.dev/).

## Available Scripts

In the project directory, you can run:

### `npm run start`

Runs the app in the development mode.\
Open [http://localhost:3000](http://localhost:3000) to view it in the browser.

The page will reload if you make edits.\
You will also see any lint errors in the console.

### `npm run build`

Builds the app for production to the `dist` folder.

It correctly bundles React in production mode and optimizes the build for the best performance.

See the section about [deployment](https://vite.dev/guide/cli.html#build) for more information.

### `npm run preview`

The command will boot up a local static web server that serves the files from `dist` folder.

It's an easy way to check if the production build looks OK in your local environment.

See the section about [deployment](https://vite.dev/guide/cli.html#build) for more information.

## Learn More

You can learn more in the [Vite documentation](https://facebook.github.io/create-react-app/docs/getting-started).

To learn React, check out the [React documentation](https://react.dev/).
