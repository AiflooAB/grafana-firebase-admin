import { DataSourcePlugin, DataQuery } from '@grafana/data';
import { DataSource } from './datasource';
import { ConfigEditor } from './ConfigEditor';
import { DataSourceOptions } from './types';

export const plugin = new DataSourcePlugin<DataSource, DataQuery, DataSourceOptions>(DataSource).setConfigEditor(
  ConfigEditor
);
