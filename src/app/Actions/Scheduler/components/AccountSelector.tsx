import React from 'react';
import {
  Button,
  FormGroup,
  Select,
  SelectOption,
  MenuToggle,
  TextInputGroup,
  TextInputGroupMain,
  TextInputGroupUtilities,
  Tooltip,
} from '@patternfly/react-core';
import { AccountResponseApi } from '@api';
import TimesIcon from '@patternfly/react-icons/dist/esm/icons/times-icon';

interface AccountTypeaheadSelectProps {
  accounts: AccountResponseApi[];
  selectedAccount: AccountResponseApi | null;
  onSelectAccount: (account: AccountResponseApi | null) => void;
  onClearAccount: () => void;
}

export const AccountTypeaheadSelect: React.FunctionComponent<AccountTypeaheadSelectProps> = ({
  accounts,
  selectedAccount,
  onSelectAccount,
  onClearAccount,
}) => {
  const [isOpen, setIsOpen] = React.useState(false);
  const [inputValue, setInputValue] = React.useState('');

  const safeAccounts = React.useMemo(() => (Array.isArray(accounts) ? accounts : []), [accounts]);

  const filteredAccounts = React.useMemo(() => {
    const q = inputValue.trim().toLowerCase();
    if (!q) return safeAccounts;

    return safeAccounts.filter(a => {
      const haystack = `${a.accountName ?? ''} ${a.accountId ?? ''}`.toLowerCase();
      return haystack.includes(q);
    });
  }, [safeAccounts, inputValue]);

  const onSelect = (_event?: React.MouseEvent<Element>, value?: string | number) => {
    const id = String(value ?? '');
    const acc = safeAccounts.find(a => a.accountId === id) ?? null;

    // Keep input in sync with selection for a predictable UX
    setInputValue(acc ? `${acc.accountName} (${acc.accountId})` : '');
    onSelectAccount(acc);
    setIsOpen(false);
  };

  return (
    <FormGroup label="Account" isRequired fieldId="account-typeahead">
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
                <div
                  className="pf-v6-u-font-size-sm pf-v6-u-font-family-mono"
                  style={{ color: 'var(--pf-t--global--text--color--subtle)' }}
                >
                  {acc.accountId}
                </div>
              </div>
            </SelectOption>
          ))
        )}
      </Select>
    </FormGroup>
  );
};
