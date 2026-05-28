/* eslint-disable */
/* tslint:disable */
// @ts-nocheck
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

import {
  ClusterEventListResponseApi,
  ClusterListResponseApi,
  ClusterRequestApi,
  ClusterResponseApi,
  GenericErrorResponseApi,
  GenericResponseApi,
  InstanceResponseApi,
  PostResponseApi,
  TagResponseApi,
} from './data-contracts';
import { ContentType, HttpClient, RequestParams } from './http-client';

export class Clusters<SecurityDataType = unknown> {
  http: HttpClient<SecurityDataType>;

  constructor(http: HttpClient<SecurityDataType>) {
    this.http = http;
  }

  /**
   * @description Paginated retrieval with optional filters.
   *
   * @tags Clusters
   * @name ClustersList
   * @summary List clusters
   * @request GET:/clusters
   */
  clustersList = (
    query?: {
      /** Account name */
      account?: string;
      /**
       * Page number
       * @default 1
       */
      page?: number;
      /**
       * Items per page
       * @default 10
       */
      page_size?: number;
      /** Cloud provider */
      provider?: string;
      /** Provider region */
      region?: string;
      /** Cluster status */
      status?: string;
    },
    params: RequestParams = {}
  ) =>
    this.http.request<ClusterListResponseApi, GenericErrorResponseApi>({
      path: `/clusters`,
      method: 'GET',
      query: query,
      type: ContentType.Json,
      format: 'json',
      ...params,
    });
  /**
   * @description Create one or multiple clusters.
   *
   * @tags Clusters
   * @name ClustersCreate
   * @summary Create clusters
   * @request POST:/clusters
   */
  clustersCreate = (clusters: ClusterRequestApi[], params: RequestParams = {}) =>
    this.http.request<PostResponseApi, GenericErrorResponseApi>({
      path: `/clusters`,
      method: 'POST',
      body: clusters,
      type: ContentType.Json,
      format: 'json',
      ...params,
    });
  /**
   * @description Return a single cluster.
   *
   * @tags Clusters
   * @name ClustersDetail
   * @summary Get cluster by ID
   * @request GET:/clusters/{id}
   */
  clustersDetail = (id: string, params: RequestParams = {}) =>
    this.http.request<ClusterResponseApi, GenericErrorResponseApi>({
      path: `/clusters/${id}`,
      method: 'GET',
      type: ContentType.Json,
      format: 'json',
      ...params,
    });
  /**
   * @description Delete a cluster by ID.
   *
   * @tags Clusters
   * @name ClustersDelete
   * @summary Delete a cluster
   * @request DELETE:/clusters/{id}
   */
  clustersDelete = (id: string, params: RequestParams = {}) =>
    this.http.request<void, GenericErrorResponseApi>({
      path: `/clusters/${id}`,
      method: 'DELETE',
      type: ContentType.Json,
      ...params,
    });
  /**
   * @description Patch mutable fields of a cluster.
   *
   * @tags Clusters
   * @name ClustersPartialUpdate
   * @summary Update a cluster
   * @request PATCH:/clusters/{id}
   */
  clustersPartialUpdate = (id: string, cluster: ClusterResponseApi, params: RequestParams = {}) =>
    this.http.request<void, void>({
      path: `/clusters/${id}`,
      method: 'PATCH',
      body: cluster,
      type: ContentType.Json,
      ...params,
    });
  /**
   * @description Paginated events for the specified cluster.
   *
   * @tags Clusters
   * @name EventsList
   * @summary List cluster events
   * @request GET:/clusters/{id}/events
   */
  eventsList = (
    id: string,
    query?: {
      /**
       * Page number
       * @default 1
       */
      page?: number;
      /**
       * Items per page
       * @default 10
       */
      page_size?: number;
    },
    params: RequestParams = {}
  ) =>
    this.http.request<ClusterEventListResponseApi, GenericErrorResponseApi>({
      path: `/clusters/${id}/events`,
      method: 'GET',
      query: query,
      type: ContentType.Json,
      format: 'json',
      ...params,
    });
  /**
   * @description Return instances belonging to the specified cluster.
   *
   * @tags Clusters
   * @name InstancesList
   * @summary Get cluster instances
   * @request GET:/clusters/{id}/instances
   */
  instancesList = (id: string, params: RequestParams = {}) =>
    this.http.request<InstanceResponseApi[], GenericErrorResponseApi>({
      path: `/clusters/${id}/instances`,
      method: 'GET',
      type: ContentType.Json,
      format: 'json',
      ...params,
    });
  /**
   * @description Request powering off a cluster.
   *
   * @tags Clusters
   * @name PowerOffCreate
   * @summary Power off a cluster
   * @request POST:/clusters/{id}/power_off
   */
  powerOffCreate = (id: string, params: RequestParams = {}) =>
    this.http.request<GenericResponseApi, GenericErrorResponseApi>({
      path: `/clusters/${id}/power_off`,
      method: 'POST',
      type: ContentType.Json,
      format: 'json',
      ...params,
    });
  /**
   * @description Request powering on a cluster.
   *
   * @tags Clusters
   * @name PowerOnCreate
   * @summary Power on a cluster
   * @request POST:/clusters/{id}/power_on
   */
  powerOnCreate = (id: string, params: RequestParams = {}) =>
    this.http.request<GenericResponseApi, GenericErrorResponseApi>({
      path: `/clusters/${id}/power_on`,
      method: 'POST',
      type: ContentType.Json,
      format: 'json',
      ...params,
    });
  /**
   * @description Return all tags for a cluster.
   *
   * @tags Clusters
   * @name TagsList
   * @summary Get cluster tags
   * @request GET:/clusters/{id}/tags
   */
  tagsList = (id: string, params: RequestParams = {}) =>
    this.http.request<TagResponseApi[], GenericErrorResponseApi>({
      path: `/clusters/${id}/tags`,
      method: 'GET',
      type: ContentType.Json,
      format: 'json',
      ...params,
    });
}
