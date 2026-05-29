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
import debounce from 'lodash.debounce';
import { ActionTypes, ActionOperations, ActionStatus } from '@app/types/types';
import { usePopperContainer } from '@app/hooks/usePopperContainer';

type AttributeMenuOption = 'Account' | 'Action' | 'Type' | 'Status' | 'Enabled';

export interface SchedulerTableToolbarProps {
  searchValue: string;
  setSearchValue: (value: string) => void;
  actionType: ActionTypes | null;
  setType: (value: ActionTypes | null) => void;
  actionOperation: ActionOperations[] | null;
  setOperation: (value: ActionOperations[] | null) => void;
  actionStatus: ActionStatus | null;
  setStatus: (value: ActionStatus | null) => void;
  actionEnabled: boolean | null;
  setEnabled: (value: boolean | null) => void;
}

export const ScheduleActionsTableToolbar: React.FunctionComponent<SchedulerTableToolbarProps> = ({
  searchValue,
  setSearchValue,
  actionType,
  setType,
  actionOperation,
  setOperation,
  actionStatus,
  setStatus,
  actionEnabled,
  setEnabled,
}) => {
  const debouncedSearch = React.useMemo(() => debounce(setSearchValue, 300), [setSearchValue]);

  React.useEffect(() => {
    return () => {
      debouncedSearch.cancel();
    };
  }, [debouncedSearch]);

  // Set up account name search input
  const searchInput = (
    <SearchInput
      placeholder="Filter by account name"
      value={searchValue}
      onChange={(_event, value) => debouncedSearch(value)}
      onClear={() => debouncedSearch('')}
    />
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
    setOperation(
      actionOperation && actionOperation.includes(selectedAction)
        ? actionOperation.filter(item => item !== selectedAction)
        : selectedAction
          ? [selectedAction, ...(actionOperation || [])]
          : []
    );
  }

  const actionToggle = (
    <MenuToggle
      ref={actionToggleRef}
      onClick={onActionMenuToggleClick}
      isExpanded={isActionMenuOpen}
      {...(actionOperation &&
        actionOperation.length > 0 && {
          badge: <Badge isRead>{actionOperation.length}</Badge>,
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
    <Menu
      ref={actionMenuRef}
      id="attribute-search-action-menu"
      onSelect={onActionMenuSelect}
      selected={actionOperation}
    >
      <MenuContent>
        <MenuList>
          <MenuItem
            hasCheckbox
            isSelected={actionOperation?.includes(ActionOperations.POWER_ON)}
            itemId={ActionOperations.POWER_ON}
          >
            {ActionOperations.POWER_ON}
          </MenuItem>
          <MenuItem
            hasCheckbox
            isSelected={actionOperation?.includes(ActionOperations.POWER_OFF)}
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

  // Type filter setup
  const [isTypeMenuOpen, setIsTypeMenuOpen] = React.useState<boolean>(false);
  const typeToggleRef = React.useRef<HTMLButtonElement>(null);
  const typeMenuRef = React.useRef<HTMLDivElement>(null);
  const { containerRef: typeContainerRef, containerElement: typeContainerElement } = usePopperContainer();

  const handleTypeMenuKeysRef = React.useRef<(event: KeyboardEvent) => void>();
  const handleTypeClickOutsideRef = React.useRef<(event: MouseEvent) => void>();

  React.useEffect(() => {
    handleTypeMenuKeysRef.current = (event: KeyboardEvent) => {
      if (isTypeMenuOpen && typeMenuRef.current?.contains(event.target as Node)) {
        if (event.key === 'Escape' || event.key === 'Tab') {
          setIsTypeMenuOpen(!isTypeMenuOpen);
          typeToggleRef.current?.focus();
        }
      }
    };

    handleTypeClickOutsideRef.current = (event: MouseEvent) => {
      if (isTypeMenuOpen && !typeMenuRef.current?.contains(event.target as Node)) {
        setIsTypeMenuOpen(false);
      }
    };
  });

  React.useEffect(() => {
    const handleKeydown = (event: KeyboardEvent) => handleTypeMenuKeysRef.current?.(event);
    const handleClick = (event: MouseEvent) => handleTypeClickOutsideRef.current?.(event);
    window.addEventListener('keydown', handleKeydown);
    window.addEventListener('click', handleClick);
    return () => {
      window.removeEventListener('keydown', handleKeydown);
      window.removeEventListener('click', handleClick);
    };
  }, [isTypeMenuOpen]);

  const onTypeMenuToggleClick = (ev: React.MouseEvent) => {
    ev.stopPropagation();
    setTimeout(() => {
      if (typeMenuRef.current) {
        const firstElement = typeMenuRef.current.querySelector('li > button:not(:disabled)');
        if (firstElement) {
          (firstElement as HTMLElement).focus();
        }
      }
    }, 0);
    setIsTypeMenuOpen(!isTypeMenuOpen);
  };

  function onTypeMenuSelect(_event: React.MouseEvent | undefined, itemId: string | number | undefined) {
    if (typeof itemId === 'undefined') {
      return;
    }

    const selectedType = itemId as ActionTypes | null;
    // Toggle type selection - clear if already selected, otherwise set new value
    setType(actionType === selectedType ? null : selectedType);
    setIsTypeMenuOpen(false);
  }

  const typeToggleLabel = (t?: ActionTypes | null) => {
    if (!t) return 'Filter by type';

    switch (t) {
      case ActionTypes.INSTANT_ACTION:
        return 'Instant Action';
      case ActionTypes.SCHEDULED_ACTION:
        return 'Scheduled Action';
      case ActionTypes.CRON_ACTION:
        return 'Cron Action';
      default:
        return 'Filter by type';
    }
  };

  const actionTypeLabel = (t: ActionTypes | null) => {
    if (!t) return '';
    if (t === ActionTypes.INSTANT_ACTION) return 'Instant Action';
    if (t === ActionTypes.SCHEDULED_ACTION) return 'Scheduled Action';
    return 'Cron Action';
  };

  const typeToggle = (
    <MenuToggle
      ref={typeToggleRef}
      onClick={onTypeMenuToggleClick}
      isExpanded={isTypeMenuOpen}
      {...(actionType && {
        badge: <Badge isRead>1</Badge>,
      })}
      style={
        {
          width: '200px',
        } as React.CSSProperties
      }
    >
      {typeToggleLabel(actionType)}
    </MenuToggle>
  );

  const typeMenu = (
    <Menu
      ref={typeMenuRef}
      id="attribute-search-type-menu"
      onSelect={onTypeMenuSelect}
      selected={actionType ? [actionType] : []}
    >
      <MenuContent>
        <MenuList>
          <MenuItem hasCheckbox isSelected={actionType === ActionTypes.INSTANT_ACTION} itemId="instant_action">
            Instant Action
          </MenuItem>
          <MenuItem hasCheckbox isSelected={actionType === ActionTypes.SCHEDULED_ACTION} itemId="scheduled_action">
            Scheduled Action
          </MenuItem>
          <MenuItem hasCheckbox isSelected={actionType === ActionTypes.CRON_ACTION} itemId="cron_action">
            Cron Action
          </MenuItem>
        </MenuList>
      </MenuContent>
    </Menu>
  );

  const typeSelect = (
    <div ref={typeContainerRef}>
      <Popper
        trigger={typeToggle}
        triggerRef={typeToggleRef}
        popper={typeMenu}
        popperRef={typeMenuRef}
        appendTo={typeContainerElement || undefined}
        isVisible={isTypeMenuOpen}
      />
    </div>
  );

  // Status filter setup
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

  const onStatusMenuToggleClick = (ev: React.MouseEvent) => {
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

  function onStatusMenuSelect(_event: React.MouseEvent | undefined, itemId: string | number | undefined) {
    if (typeof itemId === 'undefined') {
      return;
    }

    const selectedStatus = itemId as ActionStatus | null;
    // Toggle status selection
    setStatus(status === selectedStatus ? null : selectedStatus);
    setIsStatusMenuOpen(false);
  }

  const statusToggle = (
    <MenuToggle
      ref={statusToggleRef}
      onClick={onStatusMenuToggleClick}
      isExpanded={isStatusMenuOpen}
      {...(status && {
        badge: <Badge isRead>1</Badge>,
      })}
      style={
        {
          width: '200px',
        } as React.CSSProperties
      }
    >
      {status || 'Filter by status'}
    </MenuToggle>
  );

  const statusMenu = (
    <Menu
      ref={statusMenuRef}
      id="attribute-search-status-menu"
      onSelect={onStatusMenuSelect}
      selected={status ? [status] : []}
    >
      <MenuContent>
        <MenuList>
          <MenuItem hasCheckbox isSelected={status === ActionStatus.Success} itemId={ActionStatus.Success}>
            {ActionStatus.Success}
          </MenuItem>
          <MenuItem hasCheckbox isSelected={status === ActionStatus.Failed} itemId={ActionStatus.Failed}>
            {ActionStatus.Failed}
          </MenuItem>
          <MenuItem hasCheckbox isSelected={status === ActionStatus.Pending} itemId={ActionStatus.Pending}>
            {ActionStatus.Pending}
          </MenuItem>
          <MenuItem hasCheckbox isSelected={status === ActionStatus.Unknown} itemId={ActionStatus.Unknown}>
            {ActionStatus.Unknown}
          </MenuItem>
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

  // Enabled filter setup
  const [isEnabledMenuOpen, setIsEnabledMenuOpen] = React.useState<boolean>(false);
  const enabledToggleRef = React.useRef<HTMLButtonElement>(null);
  const enabledMenuRef = React.useRef<HTMLDivElement>(null);
  const { containerRef: enabledContainerRef, containerElement: enabledContainerElement } = usePopperContainer();

  const handleEnabledMenuKeysRef = React.useRef<(event: KeyboardEvent) => void>();
  const handleEnabledClickOutsideRef = React.useRef<(event: MouseEvent) => void>();

  React.useEffect(() => {
    handleEnabledMenuKeysRef.current = (event: KeyboardEvent) => {
      if (isEnabledMenuOpen && enabledMenuRef.current?.contains(event.target as Node)) {
        if (event.key === 'Escape' || event.key === 'Tab') {
          setIsEnabledMenuOpen(false);
          enabledToggleRef.current?.focus();
        }
      }
    };

    handleEnabledClickOutsideRef.current = (event: MouseEvent) => {
      if (isEnabledMenuOpen && !enabledMenuRef.current?.contains(event.target as Node)) {
        setIsEnabledMenuOpen(false);
      }
    };
  });

  React.useEffect(() => {
    const handleKeydown = (event: KeyboardEvent) => handleEnabledMenuKeysRef.current?.(event);
    const handleClick = (event: MouseEvent) => handleEnabledClickOutsideRef.current?.(event);
    window.addEventListener('keydown', handleKeydown);
    window.addEventListener('click', handleClick);
    return () => {
      window.removeEventListener('keydown', handleKeydown);
      window.removeEventListener('click', handleClick);
    };
  }, [isEnabledMenuOpen]);

  const onEnabledMenuToggleClick = (ev: React.MouseEvent) => {
    ev.stopPropagation();
    setTimeout(() => {
      if (enabledMenuRef.current) {
        const firstElement = enabledMenuRef.current.querySelector('li > button:not(:disabled)');
        if (firstElement) {
          (firstElement as HTMLElement).focus();
        }
      }
    }, 0);
    setIsEnabledMenuOpen(!isEnabledMenuOpen);
  };

  const onEnabledMenuSelect = (_event: React.MouseEvent | undefined, itemId: string | number | undefined) => {
    if (typeof itemId !== 'string') {
      return;
    }

    if (itemId === 'enabled') {
      setEnabled(true);
    } else if (itemId === 'disabled') {
      setEnabled(false);
    } else {
      setEnabled(null);
    }
  };

  const enabledToggleLabel = (v: boolean | null) => {
    if (v === true) return 'Yes';
    if (v === false) return 'No';
    return 'Filter by enabled';
  };

  const enabledToggle = (
    <MenuToggle
      ref={enabledToggleRef}
      onClick={onEnabledMenuToggleClick}
      isExpanded={isEnabledMenuOpen}
      style={{ width: '200px' } as React.CSSProperties}
    >
      {enabledToggleLabel(actionEnabled)}
    </MenuToggle>
  );

  const enabledMenu = (
    <Menu
      ref={enabledMenuRef}
      id="attribute-search-enabled-menu"
      onSelect={onEnabledMenuSelect}
      selected={actionEnabled}
    >
      <MenuContent>
        <MenuList>
          <MenuItem hasCheckbox isSelected={actionEnabled === true} itemId="enabled">
            Yes
          </MenuItem>
          <MenuItem hasCheckbox isSelected={actionEnabled === false} itemId="disabled">
            No
          </MenuItem>
        </MenuList>
      </MenuContent>
    </Menu>
  );

  const enabledSelect = (
    <div ref={enabledContainerRef}>
      <Popper
        trigger={enabledToggle}
        triggerRef={enabledToggleRef}
        popper={enabledMenu}
        popperRef={enabledMenuRef}
        appendTo={enabledContainerElement || undefined}
        isVisible={isEnabledMenuOpen}
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
          <MenuItem itemId="Type">Type</MenuItem>
          <MenuItem itemId="Status">Status</MenuItem>
          <MenuItem itemId="Enabled">Enabled</MenuItem>
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
      id="scheduler-filter-toolbar"
      clearAllFilters={() => {
        setSearchValue('');
        setOperation(null);
        setType(null);
        setStatus(null);
        setEnabled(null);
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
              labels={actionOperation || []}
              deleteLabel={(_category, chip) => onActionMenuSelect(undefined, chip as string)}
              deleteLabelGroup={() => setOperation([])}
              categoryName="Action"
              showToolbarItem={activeAttributeMenu === 'Action'}
            >
              {actionSelect}
            </ToolbarFilter>
            <ToolbarFilter
              labels={actionType ? [actionTypeLabel(actionType)] : []}
              deleteLabel={() => setType(null)}
              deleteLabelGroup={() => setType(null)}
              categoryName="Type"
              showToolbarItem={activeAttributeMenu === 'Type'}
            >
              {typeSelect}
            </ToolbarFilter>
            <ToolbarFilter
              labels={actionStatus ? [actionStatus] : []}
              deleteLabel={() => setStatus(null)}
              deleteLabelGroup={() => setStatus(null)}
              categoryName="Status"
              showToolbarItem={activeAttributeMenu === 'Status'}
            >
              {statusSelect}
            </ToolbarFilter>
            <ToolbarFilter
              labels={actionEnabled !== null ? [actionEnabled === true ? 'Enabled' : 'Disabled'] : []}
              deleteLabel={() => setEnabled(null)}
              deleteLabelGroup={() => setEnabled(null)}
              categoryName="Enabled"
              showToolbarItem={activeAttributeMenu === 'Enabled'}
            >
              {enabledSelect}
            </ToolbarFilter>
          </ToolbarGroup>
        </ToolbarToggleGroup>
      </ToolbarContent>
    </Toolbar>
  );
};

export default ScheduleActionsTableToolbar;
