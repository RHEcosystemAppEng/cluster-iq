import * as React from 'react';
import {
  Page,
  Masthead,
  MastheadToggle,
  MastheadMain,
  MastheadLogo,
  MastheadBrand,
  MastheadContent,
  PageSidebar,
  PageSidebarBody,
  PageToggleButton,
  Toolbar,
  ToolbarContent,
  ToolbarItem,
  Dropdown,
  DropdownItem,
  MenuToggle,
  ToolbarGroup,
  DropdownList,
  Tooltip,
} from '@patternfly/react-core';

import { QuestionCircleIcon, ExternalLinkAltIcon, MoonIcon, SunIcon } from '@patternfly/react-icons';
import logoImg from '../../assets/favicon.png';
import SidebarNavigation from './SidebarNavigation';
import { useUser } from '../Contexts/UserContext';
import { Link } from 'react-router-dom';
import AboutModalComponent from './AboutModal';
import { REPOSITORY_URL } from '@app/constants';
interface IAppLayout {
  children: React.ReactNode;
}

const PF_BREAKPOINT_XL = 1200;

const AppLayout: React.FunctionComponent<IAppLayout> = ({ children }) => {
  const { userEmail, setUserEmail } = useUser();
  const [isSidebarOpen, setIsSidebarOpen] = React.useState(false);
  const [isHelpMenuOpen, setIsHelpMenuOpen] = React.useState(false);
  const [isAboutModalOpen, setIsAboutModalOpen] = React.useState(false);
  const [isUserDropdownOpen, setIsUserDropdownOpen] = React.useState(false);
  const [isDarkTheme, setIsDarkTheme] = React.useState(() => {
    const saved = localStorage.getItem('theme');
    if (saved) return saved === 'dark';

    // Fallback to browser preference
    return window.matchMedia('(prefers-color-scheme: dark)').matches;
  });
  const isDesktop = () => window.innerWidth >= PF_BREAKPOINT_XL;
  const previousDesktopState = React.useRef(isDesktop());

  const defaultHelpLinks = [
    {
      label: 'Documentation',
      onClick: () => window.open(REPOSITORY_URL, '_blank'),
      isExternal: true,
    },
    {
      label: 'About',
      onClick: () => setIsAboutModalOpen(true),
    },
  ];

  const helpDropdownItems = defaultHelpLinks.map(link => {
    const content = (
      <>
        {link.label}
        {link.isExternal && (
          <span style={{ marginLeft: 'var(--pf-t--global--spacer--sm)', verticalAlign: 'middle' }}>
            {' '}
            <ExternalLinkAltIcon />
          </span>
        )}
      </>
    );
    return (
      <DropdownItem key={link.label} onClick={link.onClick} component="button">
        {content}
      </DropdownItem>
    );
  });
  // Reference https://github.com/openshift/oauth-proxy?tab=readme-ov-file#endpoint-documentation
  const handleLogout = () => {
    setUserEmail(null);
    window.location.href = '/oauth/sign_in';
  };

  const onUserDropdownToggle = () => {
    setIsUserDropdownOpen(prev => !prev);
  };

  const onUserDropdownSelect = () => {
    setIsUserDropdownOpen(false);
  };

  const userDropdownItems = [
    <DropdownItem key="logout" onClick={handleLogout}>
      Log out
    </DropdownItem>,
  ];

  const onResize = React.useCallback(() => {
    const desktop = isDesktop();
    if (desktop !== previousDesktopState.current) {
      setIsSidebarOpen(false);
    } else if (desktop) {
      setIsSidebarOpen(true);
    }
    previousDesktopState.current = desktop;
  }, []);

  const onSidebarToggle = () => {
    setIsSidebarOpen(prev => !prev);
  };

  React.useEffect(() => {
    window.addEventListener('resize', onResize);
    if (isDesktop()) {
      setIsSidebarOpen(true);
    }
    return () => {
      window.removeEventListener('resize', onResize);
    };
  }, [onResize]);

  React.useEffect(() => {
    const htmlElement = document.documentElement;
    if (isDarkTheme) {
      htmlElement.classList.add('pf-v6-theme-dark');
      localStorage.setItem('theme', 'dark');
    } else {
      htmlElement.classList.remove('pf-v6-theme-dark');
      localStorage.setItem('theme', 'light');
    }
  }, [isDarkTheme]);

  React.useEffect(() => {
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
    const handleChange = (e: MediaQueryListEvent) => {
      const saved = localStorage.getItem('theme');
      if (!saved) {
        setIsDarkTheme(e.matches);
      }
    };
    mediaQuery.addEventListener('change', handleChange);
    return () => mediaQuery.removeEventListener('change', handleChange);
  }, []);

  const toggleTheme = () => {
    setIsDarkTheme(prev => !prev);
  };

  const headerToolbar = (
    <Toolbar id="toolbar" isFullHeight isStatic className="pf-v6-u-w-100">
      <ToolbarContent className="pf-v6-u-w-100">
        <ToolbarGroup align={{ default: 'alignEnd' }}>
          <ToolbarItem>
            <Dropdown
              isOpen={isHelpMenuOpen}
              onOpenChange={setIsHelpMenuOpen}
              onSelect={() => setIsHelpMenuOpen(false)}
              popperProps={{
                position: 'right',
              }}
              toggle={toggleRef => (
                <Tooltip content="Help">
                  <MenuToggle
                    ref={toggleRef}
                    aria-label="Help dropdown"
                    variant="plain"
                    onClick={() => setIsHelpMenuOpen(!isHelpMenuOpen)}
                    isExpanded={isHelpMenuOpen}
                  >
                    <QuestionCircleIcon />
                  </MenuToggle>
                </Tooltip>
              )}
            >
              <DropdownList>{helpDropdownItems}</DropdownList>
            </Dropdown>
          </ToolbarItem>
          <ToolbarItem>
            <Tooltip content={isDarkTheme ? 'Switch to light theme' : 'Switch to dark theme'}>
              <MenuToggle
                variant="plain"
                onClick={toggleTheme}
                aria-label={isDarkTheme ? 'Switch to light theme' : 'Switch to dark theme'}
              >
                {isDarkTheme ? <SunIcon /> : <MoonIcon />}
              </MenuToggle>
            </Tooltip>
          </ToolbarItem>
          <ToolbarItem>
            <Dropdown
              isOpen={isUserDropdownOpen}
              onOpenChange={setIsUserDropdownOpen}
              onSelect={onUserDropdownSelect}
              toggle={toggleRef => (
                <MenuToggle
                  ref={toggleRef}
                  aria-label="User menu"
                  variant="plainText"
                  onClick={onUserDropdownToggle}
                  className="pf-v6-u-px-lg pf-v6-u-font-weight-normal"
                >
                  {userEmail || 'User'}
                </MenuToggle>
              )}
            >
              <DropdownList>{userDropdownItems}</DropdownList>
            </Dropdown>
          </ToolbarItem>
        </ToolbarGroup>
      </ToolbarContent>
    </Toolbar>
  );

  const header = (
    <Masthead>
      <MastheadMain>
        <MastheadToggle>
          <Tooltip content="Global navigation">
            <PageToggleButton
              isHamburgerButton
              variant="plain"
              aria-label="Global navigation"
              isSidebarOpen={isSidebarOpen}
              onSidebarToggle={onSidebarToggle}
              id="vertical-nav-toggle"
            ></PageToggleButton>
          </Tooltip>
        </MastheadToggle>
        <MastheadBrand>
          <MastheadLogo
            component={props => <Link {...props} to="/" style={{ textDecoration: 'none', color: 'inherit' }} />}
          >
            <img src={logoImg} alt="ClusterIQ" style={{ height: '36px' }} />
            <span
              style={{
                marginLeft: 'var(--pf-t--global--spacer--sm)',
                fontSize: '1.5rem',
                fontWeight: 'normal',
                verticalAlign: 'middle',
              }}
            >
              ClusterIQ
            </span>
          </MastheadLogo>
        </MastheadBrand>
      </MastheadMain>
      <MastheadContent className="pf-v6-u-w-100">{headerToolbar}</MastheadContent>
    </Masthead>
  );

  const sidebar = (
    <PageSidebar isSidebarOpen={isSidebarOpen} id="vertical-sidebar">
      <PageSidebarBody>
        <SidebarNavigation />
      </PageSidebarBody>
    </PageSidebar>
  );

  const pageId = 'primary-app-container';

  return (
    <>
      <Page masthead={header} sidebar={sidebar} mainContainerId={pageId}>
        {children}
      </Page>
      <AboutModalComponent isOpen={isAboutModalOpen} onClose={() => setIsAboutModalOpen(false)}></AboutModalComponent>
    </>
  );
};

export { AppLayout };
