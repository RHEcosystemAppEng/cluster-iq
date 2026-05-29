import React from 'react';
import {
  Button,
  FormGroup,
  Popover,
  Select,
  SelectOption,
  MenuToggle,
  TextInputGroup,
  TextInputGroupMain,
  TextInputGroupUtilities,
  Tooltip,
} from '@patternfly/react-core';
import { HelpIcon } from '@patternfly/react-icons';
import { AccountResponseApi } from '@api';
import TimesIcon from '@patternfly/react-icons/dist/esm/icons/times-icon';

export const ALL_ACCOUNTS_ID = '__all__';

interface AccountTypeaheadSelectProps {
  accounts: AccountResponseApi[];
  selectedAccount: AccountResponseApi | null;
  onSelectAccount: (account: AccountResponseApi | null) => void;
  onClearAccount: () => void;
  showAllOption?: boolean;
}

export const AccountTypeaheadSelect: React.FunctionComponent<AccountTypeaheadSelectProps> = ({
  accounts,
  selectedAccount,
  onSelectAccount,
  onClearAccount,
  showAllOption = false,
}) => {
  const [isOpen, setIsOpen] = React.useState(false);
  const [inputValue, setInputValue] = React.useState('');

  const safeAccounts = React.useMemo(() => (Array.isArray(accounts) ? accounts : []), [accounts]);

  const allAccountsEntry: AccountResponseApi = React.useMemo(
    () => ({ accountId: ALL_ACCOUNTS_ID, accountName: 'All Accounts' }),
    []
  );

  const filteredAccounts = React.useMemo(() => {
    const q = inputValue.trim().toLowerCase();
    const filtered = q
      ? safeAccounts.filter(a => {
          const haystack = `${a.accountName ?? ''} ${a.accountId ?? ''}`.toLowerCase();
          return haystack.includes(q);
        })
      : safeAccounts;

    if (showAllOption) {
      const allMatches = !q || 'all accounts'.includes(q);
      return allMatches ? [allAccountsEntry, ...filtered] : filtered;
    }
    return filtered;
  }, [safeAccounts, inputValue, showAllOption, allAccountsEntry]);

  const onSelect = (_event?: React.MouseEvent<Element>, value?: string | number) => {
    const id = String(value ?? '');
    const acc =
      (showAllOption && id === ALL_ACCOUNTS_ID ? allAccountsEntry : null) ??
      safeAccounts.find(a => a.accountId === id) ??
      null;

    setInputValue(
      acc ? (acc.accountId === ALL_ACCOUNTS_ID ? (acc.accountName ?? '') : `${acc.accountName} (${acc.accountId})`) : ''
    );
    onSelectAccount(acc);
    setIsOpen(false);
  };

  return (
    <FormGroup
      label="Account"
      isRequired
      fieldId="account-typeahead"
      labelHelp={
        <Popover
          headerContent="Account"
          bodyContent={
            showAllOption
              ? 'Select a specific account or "All Accounts" to target every account.'
              : 'The cloud provider account where the target cluster is hosted.'
          }
        >
          <Button variant="plain" aria-label="Account help" icon={<HelpIcon />} />
        </Popover>
      }
    >
      <Select
        id="account-typeahead"
        isOpen={isOpen}
        isScrollable={true}
        onOpenChange={setIsOpen}
        onSelect={onSelect}
        selected={selectedAccount?.accountId ?? null}
        toggle={toggleRef => (
          <MenuToggle
            ref={toggleRef}
            isExpanded={isOpen}
            variant="typeahead"
            onClick={() => setIsOpen(v => !v)}
            isFullWidth
          >
            <TextInputGroup isPlain>
              <TextInputGroupMain
                value={inputValue}
                onClick={() => setIsOpen(true)}
                onChange={(_e, value) => {
                  // Open menu as user types and filter client-side
                  setInputValue(value);
                  setIsOpen(true);
                }}
                placeholder="Search by name or ID..."
                aria-label="Search accounts"
              />
              <TextInputGroupUtilities {...(!inputValue ? { style: { display: 'none' } } : {})}>
                <Tooltip content="Clear input value">
                  <Button
                    icon={<TimesIcon aria-hidden />}
                    variant="plain"
                    onClick={() => {
                      setInputValue('');
                      onClearAccount();
                    }}
                    aria-label="Clear input value"
                  />
                </Tooltip>
              </TextInputGroupUtilities>
            </TextInputGroup>
          </MenuToggle>
        )}
      >
        {filteredAccounts.length === 0 ? (
          <SelectOption isDisabled value="empty">
            No results
          </SelectOption>
        ) : (
          filteredAccounts.map(acc => (
            <SelectOption key={acc.accountId} value={acc.accountId}>
              <div>
                <div>{acc.accountName}</div>
                {acc.accountId !== ALL_ACCOUNTS_ID && (
                  <div
                    className="pf-v6-u-font-size-sm pf-v6-u-font-family-mono"
                    style={{ color: 'var(--pf-t--global--text--color--subtle)' }}
                  >
                    {acc.accountId}
                  </div>
                )}
              </div>
            </SelectOption>
          ))
        )}
      </Select>
    </FormGroup>
  );
};
