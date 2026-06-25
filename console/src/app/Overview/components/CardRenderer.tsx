import React from 'react';
import { CardContentItem, CardTotalCount } from '@app/Overview/types.ts';
import { Divider, Flex, FlexItem, Stack } from '@patternfly/react-core';

export const RenderSingleIcon: React.FunctionComponent<{ content: CardContentItem[] }> = ({ content }) =>
  content[0]?.icon;

export const RenderMultiIcon: React.FunctionComponent<{ content: CardContentItem[]; totalCount?: CardTotalCount }> = ({
  content,
  totalCount,
}) => (
  <Stack hasGutter>
    <Flex display={{ default: 'inlineFlex' }} justifyContent={{ default: 'justifyContentCenter' }}>
      {content.map(({ icon, value, ref }, index) => (
        <React.Fragment key={index}>
          <Flex spaceItems={{ default: 'spaceItemsSm' }}>
            <FlexItem style={{ fontSize: '1.3em' }}>{icon}</FlexItem>
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
interface SubtitleContentItem {
  icon?: React.ReactNode;
  status: string;
  subtitle: string;
}

export const RenderWithSubtitle: React.FC<{ content: SubtitleContentItem[] }> = ({ content }) => (
  <Flex justifyContent={{ default: 'justifyContentSpaceAround' }}>
    {content.map(({ icon, status, subtitle }, index) => (
      <Flex key={index}>
        <FlexItem>{icon}</FlexItem>
        <Stack>
          <span>{status}</span>
          <span>{subtitle}</span>
        </Stack>
      </Flex>
    ))}
  </Flex>
);
