import { useState, useCallback } from 'react';

export const usePopperContainer = () => {
  const [containerElement, setContainerElement] = useState<HTMLElement | null>(null);

  const containerRef = useCallback((node: HTMLElement | null) => {
    setContainerElement(node);
  }, []);

  return { containerRef, containerElement };
};
