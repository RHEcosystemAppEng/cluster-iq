/* eslint-disable @typescript-eslint/no-explicit-any */
import React from 'react';
import { CardLayout, CardTotalCount } from '@app/Overview/types.ts';
import { RenderSingleIcon, RenderMultiIcon, RenderWithSubtitle } from '../components/CardRenderer';

export const renderContent = (content: any[], layout: CardLayout, totalCount?: CardTotalCount) => {
  switch (layout) {
    case CardLayout.SINGLE_ICON:
      return <RenderSingleIcon content={content} />;
    case CardLayout.MULTI_ICON:
      return <RenderMultiIcon content={content} totalCount={totalCount} />;
    case CardLayout.WITH_SUBTITLE:
      return <RenderWithSubtitle content={content} />;
  }
};
