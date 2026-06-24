import React from 'react';
import { Card, CardBody, CardTitle } from '@patternfly/react-core';
import { TopItemApi } from '@api';

interface TopMetricCardProps {
  title: string;
  items: TopItemApi[];
  formatValue?: (value: number) => string;
  icon?: React.ReactNode;
}

export const TopMetricCard: React.FC<TopMetricCardProps> = ({ title, items, formatValue, icon }) => {
  if (!items || items.length === 0) {
    return (
      <Card component="div" isFullHeight className="overview-card">
        <CardTitle className="pf-v6-u-text-align-center">
          {icon && <span style={{ marginRight: '0.4rem' }}>{icon}</span>}
          {title}
        </CardTitle>
        <CardBody>
          <span className="pf-v6-u-color-200">No data available</span>
        </CardBody>
      </Card>
    );
  }

  return (
    <Card component="div" isFullHeight className="overview-card">
      <CardTitle className="pf-v6-u-text-align-center">
        {icon && <span style={{ marginRight: '0.4rem' }}>{icon}</span>}
        {title}
      </CardTitle>
      <CardBody>
        <div style={{ display: 'flex', flexDirection: 'column', gap: '0.4rem' }}>
          {items.map((item, index) => (
            <div
              key={index}
              style={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                padding: '0.25rem 0',
                borderBottom:
                  index < items.length - 1 ? '1px solid var(--pf-t--global--border--color--default)' : undefined,
              }}
            >
              <span style={{ fontSize: '0.95rem' }}>{item.name || 'Unknown'}</span>
              <strong style={{ fontSize: '0.95rem' }}>
                {formatValue ? formatValue(item.clusterCount ?? 0) : (item.clusterCount ?? 0)}
              </strong>
            </div>
          ))}
        </div>
      </CardBody>
    </Card>
  );
};
