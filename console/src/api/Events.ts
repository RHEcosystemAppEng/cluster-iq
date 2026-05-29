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
  EventRequestApi,
  GenericErrorResponseApi,
  PostResponseApi,
  SystemEventListResponseApi,
} from './data-contracts';
import { ContentType, HttpClient, RequestParams } from './http-client';

export class Events<SecurityDataType = unknown> {
  http: HttpClient<SecurityDataType>;

  constructor(http: HttpClient<SecurityDataType>) {
    this.http = http;
  }

  /**
   * @description Paginated retrieval with optional filters.
   *
   * @tags Events
   * @name EventsList
   * @summary List system events
   * @request GET:/events
   */
  eventsList = (
    query?: {
      /** Action name */
      action_name?: string;
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
      /** Resource type */
      resource_type?: string;
      /** Result */
      result?: string;
      /** Severity */
      severity?: string;
      /** Triggered by */
      triggered_by?: string;
    },
    params: RequestParams = {}
  ) =>
    this.http.request<SystemEventListResponseApi, GenericErrorResponseApi>({
      path: `/events`,
      method: 'GET',
      query: query,
      type: ContentType.Json,
      format: 'json',
      ...params,
    });
  /**
   * @description Create a single event from the request body.
   *
   * @tags Events
   * @name EventsCreate
   * @summary Create event
   * @request POST:/events
   */
  eventsCreate = (event: EventRequestApi, params: RequestParams = {}) =>
    this.http.request<PostResponseApi, GenericErrorResponseApi>({
      path: `/events`,
      method: 'POST',
      body: event,
      type: ContentType.Json,
      format: 'json',
      ...params,
    });
  /**
   * @description Update an event result using the provided payload.
   *
   * @tags Events
   * @name EventsPartialUpdate
   * @summary Update event
   * @request PATCH:/events
   */
  eventsPartialUpdate = (event: EventRequestApi, params: RequestParams = {}) =>
    this.http.request<PostResponseApi, GenericErrorResponseApi>({
      path: `/events`,
      method: 'PATCH',
      body: event,
      type: ContentType.Json,
      format: 'json',
      ...params,
    });
}
