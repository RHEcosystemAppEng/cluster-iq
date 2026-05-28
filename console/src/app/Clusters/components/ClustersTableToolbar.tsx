import {
  SearchInput,
  MenuToggle,
  Menu,
  MenuContent,
  MenuList,
  MenuItem,
  Popper,
  Badge,
  Toolbar,
  ToolbarContent,
  ToolbarToggleGroup,
  ToolbarGroup,
  ToolbarItem,
  ToolbarFilter,
  Switch,
} from '@patternfly/react-core';
import { FilterIcon } from '@patternfly/react-icons';
import React from 'react';
import { ClustersTableToolbarProps } from '../types';
import debounce from 'lodash.debounce';
import { ResourceStatusApi, ProviderApi } from '@api';
import { usePopperContainer } from '@app/hooks/usePopperContainer';

export const ClustersTableToolbar: React.FunctionComponent<ClustersTableToolbarProps> = ({
  clusterNameSearch,
  setClusterNameSearch,
  accountNameSearch,
  setAccountNameSearch,
  statusSelection,
  setStatusSelection,
  providerSelections,
  setProviderSelections,
  showTerminated,
  setShowTerminated,
}) => {
  const debouncedClusterSearch = React.useMemo(() => debounce(setClusterNameSearch, 300), [setClusterNameSearch]);
  const debouncedAccountSearch = React.useMemo(() => debounce(setAccountNameSearch, 300), [setAccountNameSearch]);

  React.useEffect(() => {
    return () => {
      debouncedClusterSearch.cancel();
      debouncedAccountSearch.cancel();
    };
  }, [debouncedClusterSearch, debouncedAccountSearch]);

  const [activeAttributeMenu, setActiveAttributeMenu] = React.useState<
    'Cluster Name' | 'Account Name' | 'Status' | 'Provider'
  >('Cluster Name');

  const clusterNameInput = (
    <SearchInput
      placeholder="Filter by cluster name"
      value={clusterNameSearch}
      onChange={(_event, value) => debouncedClusterSearch(value)}
      onClear={() => debouncedClusterSearch('')}
    />
  );

  const accountNameInput = (
    <SearchInput
      placeholder="Filter by account name"
      value={accountNameSearch}
      onChange={(_event, value) => debouncedAccountSearch(value)}
      onClear={() => debouncedAccountSearch('')}
    />
  );

  // Set up status filter (only for active view)
  const [isStatusMenuOpen, setIsStatusMenuOpen] = React.useState<boolean>(false);
  const statusToggleRef = React.useRef<HTMLButtonElement>(null);
  const statusMenuRef = React.useRef<HTMLDivElement>(null);
  const { containerRef: statusContainerRef, containerElement: statusContainerElement } = usePopperContainer();

  const handleStatusMenuKeysRef = React.useRef<(event: KeyboardEvent) => void>();
  const handleStatusClickOutsideRef = React.useRef<(event: MouseEvent) => void>();

  React.useEffect(() => {
    handleStatusMenuKeysRef.current = (event: KeyboardEvent) => {
      if (isStatusMenuOpen && statusMenuRef.current?.contains(event.target as Node)) {
        if (event.key === 'Escape' || event.key === 'Tab') {
          setIsStatusMenuOpen(!isStatusMenuOpen);
          statusToggleRef.current?.focus();
        }
      }
    };

    handleStatusClickOutsideRef.current = (event: MouseEvent) => {
      if (isStatusMenuOpen && !statusMenuRef.current?.contains(event.target as Node)) {
        setIsStatusMenuOpen(false);
      }
    };
  });

  React.useEffect(() => {
    const handleKeydown = (event: KeyboardEvent) => handleStatusMenuKeysRef.current?.(event);
    const handleClick = (event: MouseEvent) => handleStatusClickOutsideRef.current?.(event);
    window.addEventListener('keydown', handleKeydown);
    window.addEventListener('click', handleClick);
    return () => {
      window.removeEventListener('keydown', handleKeydown);
      window.removeEventListener('click', handleClick);
    };
  }, [isStatusMenuOpen]);

  const onStatusToggleClick = (ev: React.MouseEvent) => {
    ev.stopPropagation();
    setTimeout(() => {
      if (statusMenuRef.current) {
        const firstElement = statusMenuRef.current.querySelector('li > button:not(:disabled)');
        if (firstElement) {
          (firstElement as HTMLElement).focus();
        }
      }
    }, 0);
    setIsStatusMenuOpen(!isStatusMenuOpen);
  };

  function onStatusSelect(_event: React.MouseEvent | undefined, itemId: string | number | undefined) {
    if (typeof itemId === 'undefined') {
      return;
    }

    setStatusSelection(itemId as ResourceStatusApi);
    setIsStatusMenuOpen(!isStatusMenuOpen);
  }

  const statusToggle = (
    <MenuToggle
      ref={statusToggleRef}
      onClick={onStatusToggleClick}
      isExpanded={isStatusMenuOpen}
      style={
        {
          width: '200px',
        } as React.CSSProperties
      }
    >
      Filter by status
    </MenuToggle>
  );

  const statusMenu = (
    <Menu ref={statusMenuRef} id="attribute-search-status-menu" onSelect={onStatusSelect} selected={statusSelection}>
      <MenuContent>
        <MenuList>
          <MenuItem itemId={ResourceStatusApi.Running}>{ResourceStatusApi.Running}</MenuItem>
          <MenuItem itemId={ResourceStatusApi.Stopped}>{ResourceStatusApi.Stopped}</MenuItem>
          <MenuItem itemId={ResourceStatusApi.Terminated}>{ResourceStatusApi.Terminated}</MenuItem>
        </MenuList>
      </MenuContent>
    </Menu>
  );

  const statusSelect = (
    <div ref={statusContainerRef}>
      <Popper
        trigger={statusToggle}
        triggerRef={statusToggleRef}
        popper={statusMenu}
        popperRef={statusMenuRef}
        appendTo={statusContainerElement || undefined}
        isVisible={isStatusMenuOpen}
      />
    </div>
  );

  // Provider filter setup
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

  const [isAttributeMenuOpen, setIsAttributeMenuOpen] = React.useState(false);
  const attributeToggleRef = React.useRef<HTMLButtonElement>(null);
  const attributeMenuRef = React.useRef<HTMLDivElement>(null);
  const { containerRef: attributeContainerRef, containerElement: attributeContainerElement } = usePopperContainer();

  const handleAttributeMenuKeysRef = React.useRef<(event: KeyboardEvent) => void>();
  const handleAttributeClickOutsideRef = React.useRef<(event: MouseEvent) => void>();

  React.useEffect(() => {
    handleAttributeMenuKeysRef.current = (event: KeyboardEvent) => {
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
    const handleKeydown = (event: KeyboardEvent) => handleAttributeMenuKeysRef.current?.(event);
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

  const onAttributeSelect = (_ev: React.MouseEvent | undefined, itemId: string | number | undefined) => {
    const selected = itemId as 'Cluster Name' | 'Account Name' | 'Status' | 'Provider';
    setActiveAttributeMenu(selected);
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
    <Menu ref={attributeMenuRef} onSelect={onAttributeSelect}>
      <MenuContent>
        <MenuList>
          <MenuItem itemId="Cluster Name">Cluster Name</MenuItem>
          <MenuItem itemId="Account Name">Account Name</MenuItem>
          <MenuItem itemId="Status">Status</MenuItem>
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
        setClusterNameSearch('');
        setAccountNameSearch('');
        setStatusSelection(null);
        setProviderSelections(null);
        setActiveAttributeMenu('Cluster Name');
      }}
    >
      <ToolbarContent>
        <ToolbarToggleGroup toggleIcon={<FilterIcon />} breakpoint="xl">
          <ToolbarGroup variant="filter-group">
            <ToolbarItem>{attributeDropdown}</ToolbarItem>
            <ToolbarFilter
              labels={clusterNameSearch !== '' ? [clusterNameSearch] : ([] as string[])}
              deleteLabel={() => setClusterNameSearch('')}
              deleteLabelGroup={() => setClusterNameSearch('')}
              categoryName="Cluster Name"
              showToolbarItem={activeAttributeMenu === 'Cluster Name'}
            >
              {clusterNameInput}
            </ToolbarFilter>
            <ToolbarFilter
              labels={accountNameSearch !== '' ? [accountNameSearch] : ([] as string[])}
              deleteLabel={() => setAccountNameSearch('')}
              deleteLabelGroup={() => setAccountNameSearch('')}
              categoryName="Account Name"
              showToolbarItem={activeAttributeMenu === 'Account Name'}
            >
              {accountNameInput}
            </ToolbarFilter>
            <ToolbarFilter
              labels={statusSelection ? [statusSelection] : []}
              deleteLabel={() => setStatusSelection(null)}
              deleteLabelGroup={() => setStatusSelection(null)}
              categoryName="Status"
              showToolbarItem={activeAttributeMenu === 'Status'}
            >
              {statusSelect}
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
        <ToolbarItem>
          <Switch
            id="show-terminated-clusters"
            label="Show terminated clusters"
            isChecked={showTerminated}
            onChange={(_event, checked) => setShowTerminated(checked)}
          />
        </ToolbarItem>
      </ToolbarContent>
    </Toolbar>
  );
};

export default ClustersTableToolbar;
