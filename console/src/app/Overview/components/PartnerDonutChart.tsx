import React from 'react';
import { Card, CardBody, CardTitle } from '@patternfly/react-core';
import { ChartDonut, ChartThemeColor } from '@patternfly/react-charts/victory';
import { TopItemApi } from '@api';

interface PartnerDonutChartProps {
  data: TopItemApi[];
}

const DONUT_COLORS = ['#06c', '#4cb140', '#009596', '#f4c145', '#ec7a08', '#7d1007', '#8481dd'];

export const PartnerDonutChart: React.FC<PartnerDonutChartProps> = ({ data }) => {
  const chartData = (data || []).map(item => ({
    x: item.name || 'Unknown',
    y: item.clusterCount ?? 0,
  }));

  const total = chartData.reduce((sum, d) => sum + d.y, 0);

  return (
    <Card component="div" isFullHeight className="overview-card">
      <CardTitle className="pf-v6-u-text-align-center">Clusters by Partner</CardTitle>
      <CardBody style={{ display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
        {chartData.length === 0 ? (
          <span className="pf-v6-u-color-200">No partner data available</span>
        ) : (
          <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', gap: '1.5rem' }}>
            <div style={{ width: '280px', height: '280px', flexShrink: 0, position: 'relative' }}>
              <ChartDonut
                constrainToVisibleArea
                data={chartData}
                themeColor={ChartThemeColor.multiOrdered}
                width={280}
                height={280}
                padding={0}
                innerRadius={90}
                style={{
                  data: {
                    fill: ({ index }: { index: number }) => DONUT_COLORS[index % DONUT_COLORS.length],
                  },
                }}
              />
              <div
                style={{
                  position: 'absolute',
                  top: '50%',
                  left: '50%',
                  transform: 'translate(-50%, -50%)',
                  textAlign: 'center',
                  color: 'var(--pf-t--global--text--color--regular)',
                  pointerEvents: 'none',
                }}
              >
                <div style={{ fontSize: '1.8rem', fontWeight: 700 }}>{total}</div>
                <div style={{ fontSize: '0.9rem' }}>Clusters</div>
              </div>
            </div>
            <div style={{ display: 'flex', flexDirection: 'column', gap: '0.35rem' }}>
              {chartData.map((d, i) => (
                <div key={i} style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                  <span
                    style={{
                      width: '12px',
                      height: '12px',
                      borderRadius: '2px',
                      backgroundColor: DONUT_COLORS[i % DONUT_COLORS.length],
                      flexShrink: 0,
                    }}
                  />
                  <span style={{ fontSize: '0.9rem', color: 'var(--pf-t--global--text--color--regular)' }}>
                    {d.x}: <strong>{d.y}</strong>
                  </span>
                </div>
              ))}
            </div>
          </div>
        )}
      </CardBody>
    </Card>
  );
};
