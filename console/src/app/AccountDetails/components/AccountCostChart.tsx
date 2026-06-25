import React, { useEffect, useMemo, useState } from 'react';
import { Card, CardBody, CardTitle, Grid, GridItem, Skeleton } from '@patternfly/react-core';
import {
  Chart,
  ChartArea,
  ChartAxis,
  ChartLine,
  ChartThemeColor,
  ChartVoronoiContainer,
} from '@patternfly/react-charts/victory';
import { api, ClusterResponseApi, DailyCostApi } from '@api';

interface AccountCostChartProps {
  accountId: string;
}

const axisTextStyle = { fill: 'var(--pf-t--global--text--color--regular)' };

// Date formatters: tickFormatDate for axis labels ("Jun 5"), formatFullDate for tooltips ("05/06/2026"), buildDateRange for card titles ("05/01/2026 — 05/06/2026")
const tickFormatDate = (t: Date) => {
  const month = t.toLocaleString('default', { month: 'short' });
  const day = t.getDate();
  return `${month} ${day}`;
};

const formatFullDate = (d: Date) =>
  `${d.getDate().toString().padStart(2, '0')}/${(d.getMonth() + 1).toString().padStart(2, '0')}/${d.getFullYear()}`;

const xAxisStyle = { tickLabels: { ...axisTextStyle, angle: -35, textAnchor: 'end' as const, fontSize: 11 } };
const yAxisStyle = { tickLabels: { ...axisTextStyle, fontSize: 11 } };

const buildDateRange = (data: { x: Date }[]): string => {
  if (data.length === 0) return 'Last 6 months';
  const first = data[0].x;
  const last = data[data.length - 1].x;
  return `${formatFullDate(first)} — ${formatFullDate(last)}`;
};

export const AccountCostChart: React.FC<AccountCostChartProps> = ({ accountId }) => {
  const [dailyCosts, setDailyCosts] = useState<DailyCostApi[]>([]);
  const [clusters, setClusters] = useState<ClusterResponseApi[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let cancelled = false;
    const fetchData = async () => {
      try {
        const [costsRes, clustersRes] = await Promise.all([
          api.accounts.dailyCostsList(accountId),
          api.accounts.clustersList(accountId),
        ]);
        if (cancelled) return;
        setDailyCosts(costsRes.data.items || []);
        setClusters(clustersRes.data.items || []);
      } catch (error) {
        if (!cancelled) console.error('Error fetching cost evolution data:', error);
      } finally {
        if (!cancelled) setLoading(false);
      }
    };

    fetchData();
    return () => {
      cancelled = true;
    };
  }, [accountId]);

  const dailyChartData = useMemo(
    () =>
      dailyCosts.map(item => ({
        x: new Date(item.date || ''),
        y: item.amount ?? 0,
      })),
    [dailyCosts]
  );

  const cumulativeChartData = useMemo(() => {
    let cumulative = 0;
    return dailyChartData.map(d => {
      cumulative += d.y;
      return { x: d.x, y: cumulative };
    });
  }, [dailyChartData]);

  const clusterCountData = useMemo(() => {
    if (dailyChartData.length === 0 || clusters.length === 0) return [];

    const sortedCreationDates = clusters
      .map(c => new Date(c.createdAt || ''))
      .filter(d => !isNaN(d.getTime()))
      .sort((a, b) => a.getTime() - b.getTime());

    return dailyChartData.map(d => {
      const count = sortedCreationDates.filter(cd => cd <= d.x).length;
      return { x: d.x, y: count };
    });
  }, [dailyChartData, clusters]);

  const dateRange = useMemo(() => buildDateRange(dailyChartData), [dailyChartData]);

  if (loading) {
    return (
      <Grid hasGutter>
        {[1, 2, 3].map(i => (
          <GridItem key={i} span={4}>
            <Card component="div" isFullHeight>
              <CardTitle className="pf-v6-u-text-align-center">Loading...</CardTitle>
              <CardBody>
                <Skeleton height="300px" />
              </CardBody>
            </Card>
          </GridItem>
        ))}
      </Grid>
    );
  }

  const maxDaily = Math.max(...dailyChartData.map(d => d.y), 0.01);
  const maxCumulative = Math.max(...cumulativeChartData.map(d => d.y), 0.01);
  const maxClusters = Math.max(...(clusterCountData.length > 0 ? clusterCountData.map(d => d.y) : [1]));

  return (
    <Grid hasGutter>
      <GridItem span={12} xl={4}>
        <Card component="div" isFullHeight>
          <CardTitle className="pf-v6-u-text-align-center">Daily Cost ({dateRange})</CardTitle>
          <CardBody>
            {dailyChartData.length === 0 ? (
              <span className="pf-v6-u-color-200">No cost data available</span>
            ) : (
              <div style={{ height: '300px', width: '100%' }}>
                <Chart
                  height={300}
                  padding={{ bottom: 80, left: 60, right: 20, top: 20 }}
                  themeColor={ChartThemeColor.blue}
                  domain={{ y: [0, maxDaily * 1.1] }}
                  scale={{ x: 'time' }}
                  containerComponent={
                    <ChartVoronoiContainer
                      labels={({ datum }: { datum: { x: Date; y: number } }) =>
                        `${formatFullDate(datum.x)}: $${datum.y.toFixed(2)}`
                      }
                      constrainToVisibleArea
                    />
                  }
                >
                  <ChartAxis fixLabelOverlap tickFormat={tickFormatDate} style={xAxisStyle} />
                  <ChartAxis dependentAxis showGrid tickFormat={(t: number) => `$${t.toFixed(2)}`} style={yAxisStyle} />
                  <ChartArea data={dailyChartData} interpolation="monotoneX" />
                </Chart>
              </div>
            )}
          </CardBody>
        </Card>
      </GridItem>

      <GridItem span={12} xl={4}>
        <Card component="div" isFullHeight>
          <CardTitle className="pf-v6-u-text-align-center">Cumulative Cost ({dateRange})</CardTitle>
          <CardBody>
            {cumulativeChartData.length === 0 ? (
              <span className="pf-v6-u-color-200">No cost data available</span>
            ) : (
              <div style={{ height: '300px', width: '100%' }}>
                <Chart
                  height={300}
                  padding={{ bottom: 80, left: 60, right: 20, top: 20 }}
                  themeColor={ChartThemeColor.blue}
                  domain={{ y: [0, maxCumulative * 1.1] }}
                  scale={{ x: 'time' }}
                  containerComponent={
                    <ChartVoronoiContainer
                      labels={({ datum }: { datum: { x: Date; y: number } }) =>
                        `${formatFullDate(datum.x)}: $${datum.y.toFixed(2)}`
                      }
                      constrainToVisibleArea
                    />
                  }
                >
                  <ChartAxis fixLabelOverlap tickFormat={tickFormatDate} style={xAxisStyle} />
                  <ChartAxis dependentAxis showGrid tickFormat={(t: number) => `$${t.toFixed(0)}`} style={yAxisStyle} />
                  <ChartArea data={cumulativeChartData} interpolation="monotoneX" />
                </Chart>
              </div>
            )}
          </CardBody>
        </Card>
      </GridItem>

      <GridItem span={12} xl={4}>
        <Card component="div" isFullHeight>
          <CardTitle className="pf-v6-u-text-align-center">Cluster Count ({dateRange})</CardTitle>
          <CardBody>
            {clusterCountData.length === 0 ? (
              <span className="pf-v6-u-color-200">No cluster data available</span>
            ) : (
              <div style={{ height: '300px', width: '100%' }}>
                <Chart
                  height={300}
                  padding={{ bottom: 80, left: 60, right: 20, top: 20 }}
                  themeColor={ChartThemeColor.green}
                  domain={{ y: [0, maxClusters + 1] }}
                  scale={{ x: 'time' }}
                  containerComponent={
                    <ChartVoronoiContainer
                      labels={({ datum }: { datum: { x: Date; y: number } }) =>
                        `${formatFullDate(datum.x)}: ${datum.y} clusters`
                      }
                      constrainToVisibleArea
                    />
                  }
                >
                  <ChartAxis fixLabelOverlap tickFormat={tickFormatDate} style={xAxisStyle} />
                  <ChartAxis dependentAxis showGrid tickFormat={(t: number) => `${Math.round(t)}`} style={yAxisStyle} />
                  <ChartLine data={clusterCountData} interpolation="stepAfter" />
                </Chart>
              </div>
            )}
          </CardBody>
        </Card>
      </GridItem>
    </Grid>
  );
};
