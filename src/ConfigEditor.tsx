import React, { ChangeEvent, PureComponent } from 'react';
import { LegacyForms } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { DataSourceOptions, SecureJsonData } from './types';

const { FormField, SecretFormField } = LegacyForms;

interface Props extends DataSourcePluginOptionsEditorProps<DataSourceOptions> {}

interface State {}

export class ConfigEditor extends PureComponent<Props, State> {
  onProjectIDChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      projectID: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  // Secure field (only sent to the backend)
  onCredentialsJSONChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({
      ...options,
      secureJsonData: {
        credentialsJSON: event.target.value,
      },
    });
  };

  onResetCredentialsJSON = () => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        credentialsJSON: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        credentialsJSON: '',
      },
    });
  };

  render() {
    const { options } = this.props;
    const { jsonData, secureJsonFields } = options;
    const secureJsonData = (options.secureJsonData || {}) as SecureJsonData;

    return (
      <div className="gf-form-group">
        <div className="gf-form">
          <FormField
            label="Project ID"
            labelWidth={12}
            inputWidth={20}
            onChange={this.onProjectIDChange}
            value={jsonData.projectID || ''}
            placeholder="Firebase project ID"
          />
        </div>
        <div className="gf-form-inline">
          <div className="gf-form">
            <SecretFormField
              isConfigured={(secureJsonFields && secureJsonFields.credentialsJSON) as boolean}
              value={secureJsonData.credentialsJSON || ''}
              label="GCP SA Credentials JSON"
              placeholder="secure json field (backend only)"
              labelWidth={12}
              inputWidth={20}
              onReset={this.onResetCredentialsJSON}
              onChange={this.onCredentialsJSONChange}
            />
          </div>
        </div>
      </div>
    );
  }
}
