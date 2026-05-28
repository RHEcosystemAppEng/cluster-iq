import { Nav, NavExpandable, NavItem, NavList } from '@patternfly/react-core';
import React from 'react';
import { NavLink, useLocation } from 'react-router-dom';

const SidebarNavigation: React.FunctionComponent = () => {
  const location = useLocation();

  const isInventoryExpanded =
    location.pathname.startsWith('/accounts') ||
    location.pathname.startsWith('/clusters') ||
    location.pathname.startsWith('/instances');

  //const isScanExpanded = location.pathname.startsWith('/scan');
  //
  // <NavExpandable title="Scan" groupId="scan-group" isExpanded={isScanExpanded} isHidden>
  // <NavItem groupId="scan-group" itemId="scan-scheduler" isActive={location.pathname === '/scan/scheduler'}>
  // <NavLink to="/scan/scheduler">Schedule</NavLink>
  // </NavItem>
  // </NavExpandable>
  //
  const isActionsExpanded = location.pathname.startsWith('/actions');

  return (
    <Nav aria-label="Nav">
      <NavList>
        <NavItem itemId="overview">
          <NavLink to="/" end>
            Overview
          </NavLink>
        </NavItem>

        <NavExpandable title="Inventory" groupId="inventory-group" isExpanded={isInventoryExpanded}>
          <NavItem groupId="inventory-group" itemId="accounts" isActive={location.pathname.startsWith('/accounts')}>
            <NavLink to="/accounts">Accounts</NavLink>
          </NavItem>
          <NavItem groupId="inventory-group" itemId="clusters" isActive={location.pathname.startsWith('/clusters')}>
            <NavLink to="/clusters">Clusters</NavLink>
          </NavItem>
          <NavItem groupId="inventory-group" itemId="instances" isActive={location.pathname.startsWith('/instances')}>
            <NavLink to="/instances">Instances</NavLink>
          </NavItem>
        </NavExpandable>

        <NavExpandable title="Actions" groupId="actions-group" isExpanded={isActionsExpanded}>
          <NavItem
            groupId="actions-group"
            itemId="actions-scheduler"
            isActive={location.pathname === '/actions/scheduler'}
          >
            <NavLink to="/actions/scheduler">Schedule</NavLink>
          </NavItem>
          <NavItem
            groupId="actions-group"
            itemId="actions-audit-logs"
            isActive={location.pathname === '/actions/audit-logs'}
          >
            <NavLink to="/actions/audit-logs">Audit Logs</NavLink>
          </NavItem>
        </NavExpandable>
      </NavList>
    </Nav>
  );
};

export default SidebarNavigation;
