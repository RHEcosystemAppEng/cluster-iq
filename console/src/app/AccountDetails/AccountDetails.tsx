import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { api, AccountResponseApi } from '@api';
import AccountsHeader from './components/AccountHeader';
import AccountsTabs from './components/AccountTabs';
import { AccountDetailsContent } from './components/AccountDetailsContent';
import { debug } from '@app/utils/debugLogs';
import { AccountClusters } from './components/AccountClusters';

const AccountDetails: React.FunctionComponent = () => {
  const { accountId } = useParams() as { accountId: string };
  const [accountData, setAccountData] = useState<AccountResponseApi | null>(null);
  const [loading, setLoading] = useState(true);
  useEffect(() => {
    const fetchData = async () => {
      try {
        debug('Fetching Account Clusters ', accountId);
        const { data: fetchedAccount } = await api.accounts.accountsDetail(accountId);
        setAccountData(fetchedAccount);
        debug('Fetched Account Clusters data:', fetchedAccount);
      } catch (error) {
        console.error('Error fetching data:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [accountId]);

  return (
    <React.Fragment>
      <AccountsHeader accountName={accountData?.accountName || accountId} label="Account" />
      <AccountsTabs
        detailsTabContent={<AccountDetailsContent loading={loading} accountData={accountData} />}
        clustersTabContent={<AccountClusters />}
      />
    </React.Fragment>
  );
};

export default AccountDetails;
