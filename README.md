# Grafana Firebase Admin Data Source Backend Plugin

This plugin adds Firebase Admin datasource that is using Firebase Admin SDK to query Firebase Auth data. Datasource simply calls `client.Users()` to fetch users data and returns user fields:
- UID
- Display name
- Email
- Phone number

For the datasource to work, you need to provide [JSON containing SA key](https://cloud.google.com/iam/docs/creating-managing-service-account-keys#creating_service_account_keys) which is needed by Firebase Admin SDK.

## What is Grafana Data Source Backend Plugin?

Grafana supports a wide range of data sources, including Prometheus, MySQL, and even Datadog. There’s a good chance you can already visualize metrics from the systems you have set up. In some cases, though, you already have an in-house metrics solution that you’d like to add to your Grafana dashboards. Grafana Data Source Plugins enables integrating such solutions with Grafana.

For more information about backend plugins, refer to the documentation on [Backend plugins](https://grafana.com/docs/grafana/latest/developers/plugins/backend/).

## Getting started

A data source backend plugin consists of both frontend and backend components.

### Frontend

1. Install dependencies

   ```bash
   yarn install
   ```

2. Build plugin in development mode or run in watch mode

   ```bash
   yarn dev
   ```

   or

   ```bash
   yarn watch
   ```

3. Build plugin in production mode

   ```bash
   yarn build
   ```

### Backend

1. Update [Grafana plugin SDK for Go](https://grafana.com/docs/grafana/latest/developers/plugins/backend/grafana-plugin-sdk-for-go/) dependency to the latest minor version:

   ```bash
   go get -u github.com/grafana/grafana-plugin-sdk-go
   go mod tidy
   ```

2. Build backend plugin binaries for Linux, Windows and Darwin:

   ```bash
   mage -v
   ```

3. List all available Mage targets for additional commands:

   ```bash
   mage -l
   ```

## Development
This plugin was developed using Node.js 14, but more recent versions should also work.

Each time you do modifications to the front-end part of the plugin you have to
execute `yarn build`, and for every change to the backend side (Go code) you have
to call `mage build:linux` (or different target if you are on a different OS)

Once everything is built, please run Grafana docker container and mount dir that
contains freshly built plugin.

```bash
docker run --rm \
   -p 3000:3000 \
   -v (pwd)/dist:/var/lib/grafana/plugins/firebase-admin \
   -e "GF_PLUGINS_ALLOW_LOADING_UNSIGNED_PLUGINS=firebase-admin" \
   --name=grafana grafana/grafana:7.4.0
```

## Sign & Pack

Grafana doesn't recommend running unsigned plugins, thats why we sign them.
It doesn't really matter which key was used to sign the plugin. Since this
component is not released very ofter, using developer key should be fine.
Otherwise a dedicated key for CI/CD pipeline should be used.

To read more about plugin signature [click here](https://grafana.com/docs/grafana/latest/developers/plugins/sign-a-plugin/)

Packing is simple
```bash
zip firebase-admin-1.0.0.zip dist -r
```

## Deployment

When starting Grafana, you can provide path to an archive containing plugin via `GF_INSTALL_PLUGINS`
https://grafana.com/docs/grafana/latest/installation/docker/#install-official-and-community-grafana-plugins

As an example:
```
GF_INSTALL_PLUGINS: "https://github.com/AiflooAB/grafana-firebase-admin/releases/download/v1.0.0/firebase-admin-1.0.0.zip;firebase-admin"
```

This plugin requires the
[firebaseauth.users.get](https://firebase.google.com/docs/projects/iam/permissions)
IAM permission in order to function.

## Learn more

- [Build a data source backend plugin tutorial](https://grafana.com/tutorials/build-a-data-source-backend-plugin)
- [Grafana documentation](https://grafana.com/docs/)
- [Grafana Tutorials](https://grafana.com/tutorials/) - Grafana Tutorials are step-by-step guides that help you make the most of Grafana
- [Grafana UI Library](https://developers.grafana.com/ui) - UI components to help you build interfaces using Grafana Design System
- [Grafana plugin SDK for Go](https://grafana.com/docs/grafana/latest/developers/plugins/backend/grafana-plugin-sdk-for-go/)
