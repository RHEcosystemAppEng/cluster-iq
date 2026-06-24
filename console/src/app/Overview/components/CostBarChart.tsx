import React from 'react';
import { Card, CardBody, CardTitle } from '@patternfly/react-core';
import { Chart, ChartBar, ChartAxis, ChartGroup, ChartThemeColor } from '@patternfly/react-charts/victory';
import { AccountCostApi } from '@api';

interface CostBarChartProps {
  data: AccountCostApi[];
}

const axisTextStyle = { fill: 'var(--pf-t--global--text--color--regular)' };

export const CostBarChart: React.FC<CostBarChartProps> = ({ data }) => {
  const chartData = (data || []).map(item => ({
    x: item.accountName || 'Unknown',
    y: item.currentMonthCost ?? 0,
  }));

  const maxCost = Math.max(...chartData.map(d => d.y), 1);

  return (
    <Card component="div" isFullHeight className="overview-card">
      <CardTitle className="pf-v6-u-text-align-center">Cost per Account (Current Month)</CardTitle>
      <CardBody>
        {chartData.length === 0 ? (
          <span className="pf-v6-u-color-200">No cost data available</span>
        ) : (
          <div style={{ height: '280px', width: '100%' }}>
            <Chart
              domainPadding={{ x: [30, 30] }}
              height={280}
              padding={{ bottom: 80, left: 80, right: 30, top: 20 }}
              themeColor={ChartThemeColor.blue}
              domain={{ y: [0, maxCost * 1.1] }}
            >
              <ChartAxis
                fixLabelOverlap
                style={{ tickLabels: { ...axisTextStyle, angle: -35, textAnchor: 'end', fontSize: 12 } }}
              />
              <ChartAxis
                dependentAxis
                showGrid
                tickFormat={(t: number) => `$${t.toFixed(0)}`}
                style={{ tickLabels: { ...axisTextStyle, fontSize: 12 } }}
              />
              <ChartGroup>
                <ChartBar data={chartData} />
              </ChartGroup>
            </Chart>
          </div>
        )}
      </CardBody>
    </Card>
  );
};
