/* eslint-disable react-refresh/only-export-components */
import * as React from 'react';

interface UserContextType {
  userEmail: string | null;
  setUserEmail: (email: string | null) => void;
}

export const UserContext = React.createContext<UserContextType>({
  userEmail: null,
  setUserEmail: () => {},
});

export const UserProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [userEmail, setUserEmail] = React.useState<string | null>(null);

  React.useEffect(() => {
    fetch(window.location.href)
      .then(response => {
        const email = response.headers.get('gap-auth');
        if (email) {
          setUserEmail(email);
          console.log('User email:', email);
        }
      })
      .catch(error => console.error('Error fetching headers:', error));
  }, []);

  return <UserContext.Provider value={{ userEmail, setUserEmail }}>{children}</UserContext.Provider>;
};

export const useUser = () => {
  const context = React.useContext(UserContext);
  if (context === undefined) {
    throw new Error('useUser must be used within a UserProvider');
  }
  return context;
};
