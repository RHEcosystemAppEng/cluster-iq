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

import { GenericErrorResponseApi, OverviewSummaryApi } from './data-contracts';
import { ContentType, HttpClient, RequestParams } from './http-client';

export class Overview<SecurityDataType = unknown> {
  http: HttpClient<SecurityDataType>;

  constructor(http: HttpClient<SecurityDataType>) {
    this.http = http;
  }

  /**
   * @description Return an aggregated summary of the inventory state.
   *
   * @tags Overview
   * @name OverviewList
   * @summary Get inventory overview
   * @request GET:/overview
   */
  overviewList = (params: RequestParams = {}) =>
    this.http.request<OverviewSummaryApi, GenericErrorResponseApi>({
      path: `/overview`,
      method: 'GET',
      type: ContentType.Json,
      format: 'json',
      ...params,
    });
}
