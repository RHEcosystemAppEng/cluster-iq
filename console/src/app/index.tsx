import * as React from 'react';
import '@patternfly/react-core/dist/styles/base.css';
import { Route, BrowserRouter as Router, Routes, useLocation } from 'react-router-dom';
import { AppLayout } from './AppLayout/AppLayout';
import Overview from './Overview/Overview';
import Clusters from './Clusters/Clusters';
import ClusterDetails from './ClusterDetails/ClusterDetails';
import AccountDetails from './AccountDetails/AccountDetails';
import NodeDetails from './NodeDetails/NodeDetails';
import AuditLogs from './Actions/AuditLogs/AuditLogs';
import Scheduler from './Actions/Scheduler/Schedule';
import Nodes from './Nodes/Nodes';
import Accounts from './Accounts/Accounts';
import { NuqsAdapter } from 'nuqs/adapters/react';
import { UserProvider } from './Contexts/UserContext';

const RouteDebugWrapper = ({ children }: { children: React.ReactNode }) => {
  const location = useLocation();

  React.useEffect(() => {
    console.log('Route changed:', {
      pathname: location.pathname,
      search: location.search,
      hash: location.hash,
    });
  }, [location]);

  return <>{children}</>;
};

const AppRoutes = (): React.ReactElement => (
  <RouteDebugWrapper>
    <Routes>
      <Route path="/" element={<Overview />} />
      <Route path="accounts" element={<Accounts />} />
      <Route path="accounts/:accountId" element={<AccountDetails />} />
      <Route path="clusters" element={<Clusters />} />
      <Route path="clusters/:clusterID" element={<ClusterDetails />} />
      <Route path="instances" element={<Nodes />} />
      <Route path="instances/:instanceID" element={<NodeDetails />} />
      <Route path="actions/scheduler" element={<Scheduler />} />
      <Route path="actions/audit-logs" element={<AuditLogs />} />
    </Routes>
  </RouteDebugWrapper>
);

const App: React.FunctionComponent = () => (
  <Router>
    <UserProvider>
      <NuqsAdapter>
        <AppLayout>
          <AppRoutes />
        </AppLayout>
      </NuqsAdapter>
    </UserProvider>
  </Router>
);

export default App;
