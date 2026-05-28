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
  GenericErrorResponseApi,
  InstanceListResponseApi,
  InstanceRequestApi,
  InstanceResponseApi,
  PostResponseApi,
} from './data-contracts';
import { ContentType, HttpClient, RequestParams } from './http-client';

export class Instances<SecurityDataType = unknown> {
  http: HttpClient<SecurityDataType>;

  constructor(http: HttpClient<SecurityDataType>) {
    this.http = http;
  }

  /**
   * @description Paginated retrieval with optional filters.
   *
   * @tags Instances
   * @name InstancesList
   * @summary List instances
   * @request GET:/instances
   */
  instancesList = (
    query?: {
      /** Cluster ID filter */
      cluster_id?: string;
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
      /** Instance status filter */
      status?: string;
    },
    params: RequestParams = {}
  ) =>
    this.http.request<InstanceListResponseApi, GenericErrorResponseApi>({
      path: `/instances`,
      method: 'GET',
      query: query,
      type: ContentType.Json,
      format: 'json',
      ...params,
    });
  /**
   * @description Create one or multiple instances.
   *
   * @tags Instances
   * @name InstancesCreate
   * @summary Create instances
   * @request POST:/instances
   */
  instancesCreate = (instances: InstanceRequestApi[], params: RequestParams = {}) =>
    this.http.request<PostResponseApi, GenericErrorResponseApi>({
      path: `/instances`,
      method: 'POST',
      body: instances,
      type: ContentType.Json,
      format: 'json',
      ...params,
    });
  /**
   * @description Return a single instance resource.
   *
   * @tags Instances
   * @name InstancesDetail
   * @summary Get instance by ID
   * @request GET:/instances/{id}
   */
  instancesDetail = (id: string, params: RequestParams = {}) =>
    this.http.request<InstanceResponseApi, GenericErrorResponseApi>({
      path: `/instances/${id}`,
      method: 'GET',
      type: ContentType.Json,
      format: 'json',
      ...params,
    });
}
