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

import { HealthCheckResponseApi } from './data-contracts';
import { ContentType, HttpClient, RequestParams } from './http-client';

export class Healthcheck<SecurityDataType = unknown> {
  http: HttpClient<SecurityDataType>;

  constructor(http: HttpClient<SecurityDataType>) {
    this.http = http;
  }

  /**
   * @description Report API process health and DB connectivity.
   *
   * @tags Health
   * @name HealthcheckList
   * @summary Health checks
   * @request GET:/healthcheck
   */
  healthcheckList = (params: RequestParams = {}) =>
    this.http.request<HealthCheckResponseApi, any>({
      path: `/healthcheck`,
      method: 'GET',
      type: ContentType.Json,
      format: 'json',
      ...params,
    });
}
