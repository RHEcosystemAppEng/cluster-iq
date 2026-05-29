import { HttpClient } from './http-client';
import { Accounts } from './Accounts';
import { Actions } from './Actions';
import { Clusters } from './Clusters';
import { Events } from './Events';
import { Instances } from './Instances';
import { Overview } from './Overview';
import { Schedule } from './Schedule';

// Initialize HTTP client with base settings
const http = new HttpClient({
  baseURL: '/api',
  timeout: 20000,
  headers: {
    'Content-Type': 'application/json',
  },
});

const PAGINATED_ENDPOINTS = ['/accounts', '/clusters', '/instances'];

http.instance.interceptors.request.use(config => {
  const isListEndpoint = PAGINATED_ENDPOINTS.some(endpoint => {
    const url = config.url || '';
    return url.includes(endpoint) && !url.match(/\/[^/]+\/[^/]+$/);
  });

  if (isListEndpoint && config.method === 'get') {
    if (!config.params?.page_size) {
      config.params = {
        page: config.params?.page || 1,
        page_size: 100,
        ...config.params,
      };
    }
  }

  return config;
});

// Initialize API classes
export const api = {
  accounts: new Accounts(http),
  actions: new Actions(http),
  clusters: new Clusters(http),
  events: new Events(http),
  instances: new Instances(http),
  overview: new Overview(http),
  schedule: new Schedule(http),
};

export const startCluster = (clusterID: string, userEmail?: string, description?: string) =>
  http.instance.post(`/clusters/${clusterID}/power_on`, {
    requester: userEmail || 'unknown',
    description: description,
  });

export const stopCluster = (clusterID: string, userEmail?: string, description?: string) =>
  http.instance.post(`/clusters/${clusterID}/power_off`, {
    requester: userEmail || 'unknown',
    description: description,
  });

export type {
  ClusterResponseApi,
  InstanceResponseApi,
  AccountResponseApi,
  ClusterListResponseApi,
  InstanceListResponseApi,
  AccountListResponseApi,
  SystemEventResponseApi,
} from './data-contracts';

export { ResourceStatusApi, ProviderApi } from './data-contracts';

export * from './data-contracts';
