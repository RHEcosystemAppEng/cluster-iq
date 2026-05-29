import {
  Badge,
  Menu,
  MenuContent,
  MenuItem,
  MenuList,
  MenuToggle,
  Popper,
  SearchInput,
  Toolbar,
  ToolbarContent,
  ToolbarFilter,
  ToolbarGroup,
  ToolbarItem,
  ToolbarToggleGroup,
} from '@patternfly/react-core';
import { FilterIcon } from '@patternfly/react-icons';
import React from 'react';
import { AuditLogsTableToolbarProps } from './types';
import debounce from 'lodash.debounce';
import { ActionOperations, ResultStatus } from '@app/types/types';
import { usePopperContainer } from '@app/hooks/usePopperContainer';
import { ProviderApi } from '@api';

type AttributeMenuOption = 'Account' | 'Provider' | 'Action' | 'Result' | 'TriggeredBy';

export const AuditLogsTableToolbar: React.FunctionComponent<AuditLogsTableToolbarProps> = ({
  searchValue,
  setSearchValue,
  action,
  setAction,
  result,
  setResult,
  triggered_by,
  setTriggeredBy,
  providerSelections,
  setProviderSelections,
}) => {
  const debouncedSearch = React.useMemo(() => debounce(setSearchValue, 300), [setSearchValue]);
  const debouncedTriggeredBy = React.useMemo(() => debounce(setTriggeredBy, 300), [setTriggeredBy]);
  const debouncedResult = React.useMemo(() => debounce(setResult, 300), [setResult]);

  React.useEffect(() => {
    return () => {
      debouncedSearch.cancel();
      debouncedTriggeredBy.cancel();
      debouncedResult.cancel();
    };
  }, [debouncedSearch, debouncedTriggeredBy, debouncedResult]);

  // Set up name search input
  const searchInput = (
    <SearchInput
      placeholder="Filter by account name"
      value={searchValue}
      onChange={(_event, value) => debouncedSearch(value)}
      onClear={() => debouncedSearch('')}
    />
  );
  // Set up triggered by search input
  const triggeredByInput = (
    <SearchInput
      placeholder="Filter by user"
      value={searchValue}
      onChange={(_event, value) => debouncedTriggeredBy(value)}
      onClear={() => debouncedTriggeredBy('')}
    />
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

  // Actions filter setup
  const [isActionMenuOpen, setIsActionMenuOpen] = React.useState<boolean>(false);
  const actionToggleRef = React.useRef<HTMLButtonElement>(null);
  const actionMenuRef = React.useRef<HTMLDivElement>(null);
  const { containerRef: actionContainerRef, containerElement: actionContainerElement } = usePopperContainer();

  const handleActionMenuKeysRef = React.useRef<(event: KeyboardEvent) => void>();
  const handleActionClickOutsideRef = React.useRef<(event: MouseEvent) => void>();

  React.useEffect(() => {
    handleActionMenuKeysRef.current = (event: KeyboardEvent) => {
      if (isActionMenuOpen && actionMenuRef.current?.contains(event.target as Node)) {
        if (event.key === 'Escape' || event.key === 'Tab') {
          setIsActionMenuOpen(!isActionMenuOpen);
          actionToggleRef.current?.focus();
        }
      }
    };

    handleActionClickOutsideRef.current = (event: MouseEvent) => {
      if (isActionMenuOpen && !actionMenuRef.current?.contains(event.target as Node)) {
        setIsActionMenuOpen(false);
      }
    };
  });

  React.useEffect(() => {
    const handleKeydown = (event: KeyboardEvent) => handleActionMenuKeysRef.current?.(event);
    const handleClick = (event: MouseEvent) => handleActionClickOutsideRef.current?.(event);
    window.addEventListener('keydown', handleKeydown);
    window.addEventListener('click', handleClick);
    return () => {
      window.removeEventListener('keydown', handleKeydown);
      window.removeEventListener('click', handleClick);
    };
  }, [isActionMenuOpen]);

  const onActionMenuToggleClick = (ev: React.MouseEvent) => {
    ev.stopPropagation();
    setTimeout(() => {
      if (actionMenuRef.current) {
        const firstElement = actionMenuRef.current.querySelector('li > button:not(:disabled)');
        if (firstElement) {
          (firstElement as HTMLElement).focus();
        }
      }
    }, 0);
    setIsActionMenuOpen(!isActionMenuOpen);
  };

  function onActionMenuSelect(_event: React.MouseEvent | undefined, itemId: string | number | undefined) {
    if (typeof itemId === 'undefined') {
      return;
    }

    const selectedAction = itemId as ActionOperations;
    setAction(
      action && action.includes(selectedAction)
        ? action.filter(item => item !== selectedAction)
        : selectedAction
          ? [selectedAction, ...(action || [])]
          : []
    );
  }

  const actionToggle = (
    <MenuToggle
      ref={actionToggleRef}
      onClick={onActionMenuToggleClick}
      isExpanded={isActionMenuOpen}
      {...(action &&
        action.length > 0 && {
          badge: <Badge isRead>{action.length}</Badge>,
        })}
      style={
        {
          width: '200px',
        } as React.CSSProperties
      }
    >
      Filter by action
    </MenuToggle>
  );

  const actionMenu = (
    <Menu ref={actionMenuRef} id="attribute-search-action-menu" onSelect={onActionMenuSelect} selected={action}>
      <MenuContent>
        <MenuList>
          <MenuItem
            hasCheckbox
            isSelected={action?.includes(ActionOperations.POWER_ON)}
            itemId={ActionOperations.POWER_ON}
          >
            {ActionOperations.POWER_ON}
          </MenuItem>
          <MenuItem
            hasCheckbox
            isSelected={action?.includes(ActionOperations.POWER_OFF)}
            itemId={ActionOperations.POWER_OFF}
          >
            {ActionOperations.POWER_OFF}
          </MenuItem>
        </MenuList>
      </MenuContent>
    </Menu>
  );

  const actionSelect = (
    <div ref={actionContainerRef}>
      <Popper
        trigger={actionToggle}
        triggerRef={actionToggleRef}
        popper={actionMenu}
        popperRef={actionMenuRef}
        appendTo={actionContainerElement || undefined}
        isVisible={isActionMenuOpen}
      />
    </div>
  );

  // Result filter setup
  const [isResultMenuOpen, setIsResultMenuOpen] = React.useState<boolean>(false);
  const resultToggleRef = React.useRef<HTMLButtonElement>(null);
  const resultMenuRef = React.useRef<HTMLDivElement>(null);
  const { containerRef: resultContainerRef, containerElement: resultContainerElement } = usePopperContainer();

  const handleResultMenuKeysRef = React.useRef<(event: KeyboardEvent) => void>();
  const handleResultClickOutsideRef = React.useRef<(event: MouseEvent) => void>();

  React.useEffect(() => {
    handleResultMenuKeysRef.current = (event: KeyboardEvent) => {
      if (isResultMenuOpen && resultMenuRef.current?.contains(event.target as Node)) {
        if (event.key === 'Escape' || event.key === 'Tab') {
          setIsResultMenuOpen(!isResultMenuOpen);
          resultToggleRef.current?.focus();
        }
      }
    };

    handleResultClickOutsideRef.current = (event: MouseEvent) => {
      if (isResultMenuOpen && !resultMenuRef.current?.contains(event.target as Node)) {
        setIsResultMenuOpen(false);
      }
    };
  });

  React.useEffect(() => {
    const handleKeydown = (event: KeyboardEvent) => handleResultMenuKeysRef.current?.(event);
    const handleClick = (event: MouseEvent) => handleResultClickOutsideRef.current?.(event);
    window.addEventListener('keydown', handleKeydown);
    window.addEventListener('click', handleClick);
    return () => {
      window.removeEventListener('keydown', handleKeydown);
      window.removeEventListener('click', handleClick);
    };
  }, [isResultMenuOpen]);

  const onResultMenuToggleClick = (ev: React.MouseEvent) => {
    ev.stopPropagation();
    setTimeout(() => {
      if (resultMenuRef.current) {
        const firstElement = resultMenuRef.current.querySelector('li > button:not(:disabled)');
        if (firstElement) {
          (firstElement as HTMLElement).focus();
        }
      }
    }, 0);
    setIsResultMenuOpen(!isResultMenuOpen);
  };

  function onResultMenuSelect(_event: React.MouseEvent | undefined, itemId: string | number | undefined) {
    if (typeof itemId === 'undefined') {
      return;
    }

    const selectedResult = itemId as ResultStatus;
    setResult(
      result && result.includes(selectedResult)
        ? result.filter(item => item !== selectedResult)
        : selectedResult
          ? [selectedResult, ...(result || [])]
          : []
    );
  }

  const resultToggle = (
    <MenuToggle
      ref={resultToggleRef}
      onClick={onResultMenuToggleClick}
      isExpanded={isResultMenuOpen}
      {...(result &&
        result.length > 0 && {
          badge: <Badge isRead>{result.length}</Badge>,
        })}
      style={
        {
          width: '200px',
        } as React.CSSProperties
      }
    >
      Filter by result
    </MenuToggle>
  );

  const resultMenu = (
    <Menu ref={resultMenuRef} id="attribute-search-result-menu" onSelect={onResultMenuSelect} selected={result}>
      <MenuContent>
        <MenuList>
          <MenuItem hasCheckbox isSelected={result?.includes(ResultStatus.Success)} itemId={ResultStatus.Success}>
            {ResultStatus.Success}
          </MenuItem>
          <MenuItem hasCheckbox isSelected={result?.includes(ResultStatus.Failed)} itemId={ResultStatus.Failed}>
            {ResultStatus.Failed}
          </MenuItem>
          <MenuItem hasCheckbox isSelected={result?.includes(ResultStatus.Warning)} itemId={ResultStatus.Warning}>
            {ResultStatus.Warning}
          </MenuItem>
        </MenuList>
      </MenuContent>
    </Menu>
  );

  const resultSelect = (
    <div ref={resultContainerRef}>
      <Popper
        trigger={resultToggle}
        triggerRef={resultToggleRef}
        popper={resultMenu}
        popperRef={resultMenuRef}
        appendTo={resultContainerElement || undefined}
        isVisible={isResultMenuOpen}
      />
    </div>
  );

  // Attribute selector setup
  const [activeAttributeMenu, setActiveAttributeMenu] = React.useState<AttributeMenuOption>('Account');
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

  // Only show Status option in attribute menu for active view
  const attributeMenu = (
    <Menu
      ref={attributeMenuRef}
      onSelect={(_ev, itemId) => {
        const selected = itemId?.toString() as AttributeMenuOption;
        setActiveAttributeMenu(selected);
        setIsAttributeMenuOpen(!isAttributeMenuOpen);
      }}
    >
      <MenuContent>
        <MenuList>
          <MenuItem itemId="Account">Account</MenuItem>
          <MenuItem itemId="Action">Action</MenuItem>
          <MenuItem itemId="Result">Result</MenuItem>
          <MenuItem itemId="Provider">Provider</MenuItem>
          <MenuItem itemId="TriggeredBy">Triggered by</MenuItem>
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
        setProviderSelections(null);
        setAction(null);
        setResult(null);
        setTriggeredBy('');
      }}
    >
      <ToolbarContent>
        <ToolbarToggleGroup toggleIcon={<FilterIcon />} breakpoint="xl">
          <ToolbarGroup variant="filter-group">
            <ToolbarItem>{attributeDropdown}</ToolbarItem>
            <ToolbarFilter
              labels={searchValue !== '' ? [searchValue] : []}
              deleteLabel={() => setSearchValue('')}
              deleteLabelGroup={() => setSearchValue('')}
              categoryName="Account"
              showToolbarItem={activeAttributeMenu === 'Account'}
            >
              {searchInput}
            </ToolbarFilter>
            <ToolbarFilter
              labels={triggered_by !== '' ? [triggered_by] : []}
              deleteLabel={() => setTriggeredBy('')}
              deleteLabelGroup={() => setTriggeredBy('')}
              categoryName="TriggeredBy"
              showToolbarItem={activeAttributeMenu === 'TriggeredBy'}
            >
              {triggeredByInput}
            </ToolbarFilter>
            <ToolbarFilter
              labels={result || []}
              deleteLabel={(_category, chip) => onResultMenuSelect(undefined, chip as string)}
              deleteLabelGroup={() => setResult([])}
              categoryName="Result"
              showToolbarItem={activeAttributeMenu === 'Result'}
            >
              {resultSelect}
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
            <ToolbarFilter
              labels={action || []}
              deleteLabel={(_category, chip) => onActionMenuSelect(undefined, chip as string)}
              deleteLabelGroup={() => setAction([])}
              categoryName="Action"
              showToolbarItem={activeAttributeMenu === 'Action'}
            >
              {actionSelect}
            </ToolbarFilter>
          </ToolbarGroup>
        </ToolbarToggleGroup>
      </ToolbarContent>
    </Toolbar>
  );
};

export default AuditLogsTableToolbar;
