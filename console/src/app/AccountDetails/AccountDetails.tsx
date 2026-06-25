import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { api, AccountResponseApi } from '@api';
import AccountsHeader from './components/AccountHeader';
import AccountsTabs from './components/AccountTabs';
import { AccountDetailsContent } from './components/AccountDetailsContent';
import { AccountCostChart } from './components/AccountCostChart';
import { debug } from '@app/utils/debugLogs';
import { AccountClusters } from './components/AccountClusters';
import { useDocumentTitle } from '@app/utils/useDocumentTitle';

const AccountDetails: React.FunctionComponent = () => {
  const { accountId } = useParams() as { accountId: string };
  const [accountData, setAccountData] = useState<AccountResponseApi | null>(null);
  useDocumentTitle(`${accountData?.accountName || accountId} — ClusterIQ`);
  const [loading, setLoading] = useState(true);
  useEffect(() => {
    let cancelled = false;
    const fetchData = async () => {
      try {
        debug('Fetching Account Clusters ', accountId);
        const { data: fetchedAccount } = await api.accounts.accountsDetail(accountId);
        if (cancelled) return;
        setAccountData(fetchedAccount);
        debug('Fetched Account Clusters data:', fetchedAccount);
      } catch (error) {
        if (!cancelled) console.error('Error fetching data:', error);
      } finally {
        if (!cancelled) setLoading(false);
      }
    };

    fetchData();
    return () => {
      cancelled = true;
    };
  }, [accountId]);

  return (
    <React.Fragment>
      <AccountsHeader accountName={accountData?.accountName || accountId} accountId={accountId} />
      <AccountsTabs
        detailsTabContent={<AccountDetailsContent loading={loading} accountData={accountData} />}
        clustersTabContent={<AccountClusters />}
        costsTabContent={<AccountCostChart accountId={accountId} />}
      />
    </React.Fragment>
  );
};

export default AccountDetails;
