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

import { ActionListResponseApi, GenericErrorResponseApi, ScheduledActionApi } from './data-contracts';
import { ContentType, HttpClient, RequestParams } from './http-client';

export class Schedule<SecurityDataType = unknown> {
  http: HttpClient<SecurityDataType>;

  constructor(http: HttpClient<SecurityDataType>) {
    this.http = http;
  }

  /**
   * @description Paginated retrieval with optional filters.
   *
   * @tags Actions
   * @name ScheduleList
   * @summary List scheduled actions
   * @request GET:/actions
   */
  scheduleList = (
    query?: {
      /** Enabled state filter (true/false) */
      enabled?: string;
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
      /** Status filter */
      status?: string;
    },
    params: RequestParams = {}
  ) =>
    this.http.request<ActionListResponseApi, GenericErrorResponseApi>({
      path: `/actions`,
      method: 'GET',
      query: query,
      type: ContentType.Json,
      format: 'json',
      ...params,
    });
  /**
   * @description Return a scheduled action resource.
   *
   * @tags Actions
   * @name ScheduleDetail
   * @summary Get scheduled action by ID
   * @request GET:/actions/{id}
   */
  scheduleDetail = (id: string, params: RequestParams = {}) =>
    this.http.request<ScheduledActionApi, GenericErrorResponseApi>({
      path: `/actions/${id}`,
      method: 'GET',
      type: ContentType.Json,
      format: 'json',
      ...params,
    });
  /**
   * @description Disable a scheduled action.
   *
   * @tags Actions
   * @name DisablePartialUpdate
   * @summary Disable scheduled action
   * @request PATCH:/actions/{id}/disable
   */
  disablePartialUpdate = (id: string, params: RequestParams = {}) =>
    this.http.request<void, GenericErrorResponseApi>({
      path: `/actions/${id}/disable`,
      method: 'PATCH',
      ...params,
    });
  /**
   * @description Enable a scheduled action.
   *
   * @tags Actions
   * @name EnablePartialUpdate
   * @summary Enable scheduled action
   * @request PATCH:/actions/{id}/enable
   */
  enablePartialUpdate = (id: string, params: RequestParams = {}) =>
    this.http.request<void, GenericErrorResponseApi>({
      path: `/actions/${id}/enable`,
      method: 'PATCH',
      ...params,
    });
}
