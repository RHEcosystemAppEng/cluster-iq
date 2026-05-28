import {
  SearchInput,
  MenuToggle,
  Badge,
  Menu,
  MenuContent,
  MenuList,
  MenuItem,
  Popper,
  Toolbar,
  ToolbarContent,
  ToolbarToggleGroup,
  ToolbarGroup,
  ToolbarItem,
  ToolbarFilter,
} from '@patternfly/react-core';
import { FilterIcon } from '@patternfly/react-icons';
import React from 'react';
import { AccountsToolbarProps } from '../types';
import { ProviderApi } from '@api';
import debounce from 'lodash.debounce';
import { usePopperContainer } from '@app/hooks/usePopperContainer';

export const AccountsToolbar: React.FunctionComponent<AccountsToolbarProps> = ({
  searchValue,
  setSearchValue,
  providerSelections,
  setProviderSelections,
}) => {
  const debouncedSearch = React.useMemo(() => debounce(setSearchValue, 300), [setSearchValue]);

  React.useEffect(() => {
    return () => {
      debouncedSearch.cancel();
    };
  }, [debouncedSearch]);

  const searchInput = (
    <SearchInput
      placeholder="Filter by account name"
      value={searchValue}
      onChange={(_event, value) => debouncedSearch(value)}
      onClear={() => debouncedSearch('')}
    />
  );

  const [isProviderMenuOpen, setIsProviderMenuOpen] = React.useState<boolean>(false);
  const providerToggleRef = React.useRef<HTMLButtonElement>(null);
  const providerMenuRef = React.useRef<HTMLDivElement>(null);
  const { containerRef: providerContainerRef, containerElement: providerContainerElement } = usePopperContainer();

  const handleProviderMenuKeysRef = React.useRef<(event: KeyboardEvent) => void>();
  const handleProviderClickOutsideRef = React.useRef<(event: MouseEvent) => void>();

  React.useEffect(() => {
    handleProviderMenuKeysRef.current = (event: KeyboardEvent) => {
      if (isProviderMenuOpen && providerMenuRef.current?.contains(event.target as Node)) {
        if (event.key === 'Escape' || event.key === 'Tab') {
          setIsProviderMenuOpen(!isProviderMenuOpen);
          providerToggleRef.current?.focus();
        }
      }
    };

    handleProviderClickOutsideRef.current = (event: MouseEvent) => {
      if (isProviderMenuOpen && !providerMenuRef.current?.contains(event.target as Node)) {
        setIsProviderMenuOpen(false);
      }
    };
  });

  React.useEffect(() => {
    const handleKeydown = (event: KeyboardEvent) => handleProviderMenuKeysRef.current?.(event);
    const handleClick = (event: MouseEvent) => handleProviderClickOutsideRef.current?.(event);
    window.addEventListener('keydown', handleKeydown);
    window.addEventListener('click', handleClick);
    return () => {
      window.removeEventListener('keydown', handleKeydown);
      window.removeEventListener('click', handleClick);
    };
  }, [isProviderMenuOpen]);

  const onProviderMenuToggleClick = (ev: React.MouseEvent) => {
    ev.stopPropagation();
    setTimeout(() => {
      if (providerMenuRef.current) {
        const firstElement = providerMenuRef.current.querySelector('li > button:not(:disabled)');
        if (firstElement) {
          (firstElement as HTMLElement).focus();
        }
      }
    }, 0);
    setIsProviderMenuOpen(!isProviderMenuOpen);
  };

  function onProviderMenuSelect(_event: React.MouseEvent | undefined, itemId: string | number | undefined) {
    if (typeof itemId === 'undefined') {
      return;
    }

    const provider = itemId as ProviderApi;
    setProviderSelections(
      providerSelections && providerSelections.includes(provider)
        ? providerSelections.filter(selection => selection !== provider)
        : provider
          ? [provider, ...(providerSelections || [])]
          : []
    );
  }

  const providerToggle = (
    <MenuToggle
      ref={providerToggleRef}
      onClick={onProviderMenuToggleClick}
      isExpanded={isProviderMenuOpen}
      {...(providerSelections &&
        providerSelections.length > 0 && {
          badge: <Badge isRead>{providerSelections.length}</Badge>,
        })}
      style={
        {
          width: '200px',
        } as React.CSSProperties
      }
    >
      Filter by provider
    </MenuToggle>
  );

  const providerMenu = (
    <Menu
      ref={providerMenuRef}
      id="attribute-search-provider-menu"
      onSelect={onProviderMenuSelect}
      selected={providerSelections}
    >
      <MenuContent>
        <MenuList>
          <MenuItem
            hasCheckbox
            isSelected={providerSelections?.includes(ProviderApi.AWSProvider)}
            itemId={ProviderApi.AWSProvider}
          >
            AWS
          </MenuItem>
          <MenuItem
            hasCheckbox
            isSelected={providerSelections?.includes(ProviderApi.GCPProvider)}
            itemId={ProviderApi.GCPProvider}
          >
            Google Cloud
          </MenuItem>
          <MenuItem
            hasCheckbox
            isSelected={providerSelections?.includes(ProviderApi.AzureProvider)}
            itemId={ProviderApi.AzureProvider}
          >
            Azure
          </MenuItem>
        </MenuList>
      </MenuContent>
    </Menu>
  );

  const providerSelect = (
    <div ref={providerContainerRef}>
      <Popper
        trigger={providerToggle}
        triggerRef={providerToggleRef}
        popper={providerMenu}
        popperRef={providerMenuRef}
        appendTo={providerContainerElement || undefined}
        isVisible={isProviderMenuOpen}
      />
    </div>
  );

  // Set up attribute selector
  const [activeAttributeMenu, setActiveAttributeMenu] = React.useState<'Account' | 'Provider'>('Account');
  const [isAttributeMenuOpen, setIsAttributeMenuOpen] = React.useState(false);
  const attributeToggleRef = React.useRef<HTMLButtonElement>(null);
  const attributeMenuRef = React.useRef<HTMLDivElement>(null);
  const { containerRef: attributeContainerRef, containerElement: attributeContainerElement } = usePopperContainer();

  const handleAttribueMenuKeysRef = React.useRef<(event: KeyboardEvent) => void>();
  const handleAttributeClickOutsideRef = React.useRef<(event: MouseEvent) => void>();

  React.useEffect(() => {
    handleAttribueMenuKeysRef.current = (event: KeyboardEvent) => {
      if (!isAttributeMenuOpen) {
        return;
      }
      if (
        attributeMenuRef.current?.contains(event.target as Node) ||
        attributeToggleRef.current?.contains(event.target as Node)
      ) {
        if (event.key === 'Escape' || event.key === 'Tab') {
          setIsAttributeMenuOpen(!isAttributeMenuOpen);
          attributeToggleRef.current?.focus();
        }
      }
    };

    handleAttributeClickOutsideRef.current = (event: MouseEvent) => {
      if (isAttributeMenuOpen && !attributeMenuRef.current?.contains(event.target as Node)) {
        setIsAttributeMenuOpen(false);
      }
    };
  });

  React.useEffect(() => {
    const handleKeydown = (event: KeyboardEvent) => handleAttribueMenuKeysRef.current?.(event);
    const handleClick = (event: MouseEvent) => handleAttributeClickOutsideRef.current?.(event);
    window.addEventListener('keydown', handleKeydown);
    window.addEventListener('click', handleClick);
    return () => {
      window.removeEventListener('keydown', handleKeydown);
      window.removeEventListener('click', handleClick);
    };
  }, [isAttributeMenuOpen]);

  const onAttributeToggleClick = (ev: React.MouseEvent) => {
    ev.stopPropagation();
    setTimeout(() => {
      if (attributeMenuRef.current) {
        const firstElement = attributeMenuRef.current.querySelector('li > button:not(:disabled)');
        if (firstElement) {
          (firstElement as HTMLElement).focus();
        }
      }
    }, 0);
    setIsAttributeMenuOpen(!isAttributeMenuOpen);
  };

  const attributeToggle = (
    <MenuToggle
      ref={attributeToggleRef}
      onClick={onAttributeToggleClick}
      isExpanded={isAttributeMenuOpen}
      icon={<FilterIcon />}
    >
      {activeAttributeMenu}
    </MenuToggle>
  );
  const attributeMenu = (
    <Menu
      ref={attributeMenuRef}
      onSelect={(_ev, itemId) => {
        setActiveAttributeMenu(itemId?.toString() as 'Account' | 'Provider');
        setIsAttributeMenuOpen(!isAttributeMenuOpen);
      }}
    >
      <MenuContent>
        <MenuList>
          <MenuItem itemId="Account">Account</MenuItem>
          <MenuItem itemId="Provider">Provider</MenuItem>
        </MenuList>
      </MenuContent>
    </Menu>
  );

  const attributeDropdown = (
    <div ref={attributeContainerRef}>
      <Popper
        trigger={attributeToggle}
        triggerRef={attributeToggleRef}
        popper={attributeMenu}
        popperRef={attributeMenuRef}
        appendTo={attributeContainerElement || undefined}
        isVisible={isAttributeMenuOpen}
      />
    </div>
  );

  return (
    <Toolbar
      id="attribute-search-filter-toolbar"
      clearAllFilters={() => {
        setSearchValue('');
        setProviderSelections([]);
      }}
    >
      <ToolbarContent>
        <ToolbarToggleGroup toggleIcon={<FilterIcon />} breakpoint="xl">
          <ToolbarGroup variant="filter-group">
            <ToolbarItem>{attributeDropdown}</ToolbarItem>
            <ToolbarFilter
              labels={searchValue !== '' ? [searchValue] : ([] as string[])}
              deleteLabel={() => setSearchValue('')}
              deleteLabelGroup={() => setSearchValue('')}
              categoryName="Name"
              showToolbarItem={activeAttributeMenu === 'Account'}
            >
              {searchInput}
            </ToolbarFilter>

            <ToolbarFilter
              labels={providerSelections || []}
              deleteLabel={(_category, chip) => onProviderMenuSelect(undefined, chip as string)}
              deleteLabelGroup={() => setProviderSelections([])}
              categoryName="Provider"
              showToolbarItem={activeAttributeMenu === 'Provider'}
            >
              {providerSelect}
            </ToolbarFilter>
          </ToolbarGroup>
        </ToolbarToggleGroup>
      </ToolbarContent>
    </Toolbar>
  );
};

export default AccountsToolbar;
