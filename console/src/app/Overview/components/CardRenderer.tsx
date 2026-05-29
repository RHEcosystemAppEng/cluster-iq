/* eslint-disable @typescript-eslint/no-explicit-any */
import React from 'react';
import { CardTotalCount } from '@app/Overview/types.ts';
import { Divider, Flex, FlexItem, Stack } from '@patternfly/react-core';

// TODO Avoid any
export const RenderSingleIcon: React.FunctionComponent<{ content: any[] }> = ({ content }) => content[0]?.icon;
// TODO Avoid any
export const RenderMultiIcon: React.FunctionComponent<{ content: any[]; totalCount?: CardTotalCount }> = ({
  content,
  totalCount,
}) => (
  <Stack hasGutter>
    <Flex display={{ default: 'inlineFlex' }} justifyContent={{ default: 'justifyContentCenter' }}>
      {content.map(({ icon, value, ref }, index) => (
        <React.Fragment key={index}>
          <Flex spaceItems={{ default: 'spaceItemsSm' }}>
            <FlexItem>{icon}</FlexItem>
            <FlexItem>{ref ? <a href={ref}>{value}</a> : <span>{value}</span>}</FlexItem>
          </Flex>
          {content.length > 1 && index < content.length - 1 && <Divider orientation={{ default: 'vertical' }} />}
        </React.Fragment>
      ))}
    </Flex>
    {totalCount && (
      <Flex justifyContent={{ default: 'justifyContentCenter' }} spaceItems={{ default: 'spaceItemsSm' }}>
        <FlexItem>{totalCount.icon}</FlexItem>
        <FlexItem>
          <span style={{ fontWeight: 500 }}>
            {totalCount.label}: {totalCount.value}
          </span>
        </FlexItem>
      </Flex>
    )}
  </Stack>
);
// TODO Avoid any
export const RenderWithSubtitle: React.FC<{ content: any[] }> = ({ content }) => (
  <Flex justifyContent={{ default: 'justifyContentSpaceAround' }}>
    {content.map(({ icon, status, subtitle }, index) => (
      <Flex key={index}>
        <FlexItem>{icon}</FlexItem>
        <Stack>
          <a href="#">{status}</a>
          <span>{subtitle}</span>
        </Stack>
      </Flex>
    ))}
  </Flex>
);
