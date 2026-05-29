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

import { ActionRequestApi, GenericErrorResponseApi, PostResponseApi } from './data-contracts';
import { ContentType, HttpClient, RequestParams } from './http-client';

export class Actions<SecurityDataType = unknown> {
  http: HttpClient<SecurityDataType>;

  constructor(http: HttpClient<SecurityDataType>) {
    this.http = http;
  }

  /**
   * @description Create one or multiple actions.
   *
   * @tags Actions
   * @name ActionsCreate
   * @summary Create actions
   * @request POST:/actions
   */
  actionsCreate = (actions: ActionRequestApi[], params: RequestParams = {}) =>
    this.http.request<PostResponseApi, GenericErrorResponseApi>({
      path: `/actions`,
      method: 'POST',
      body: actions,
      type: ContentType.Json,
      format: 'json',
      ...params,
    });
  /**
   * @description Enables one or multiple actions.
   *
   * @tags Actions
   * @name ActionsEnable
   * @summary Enable actions
   * @request PATCH:/actions
   */
  actionsEnable = (id: string, params: RequestParams = {}) =>
    this.http.request<PostResponseApi, GenericErrorResponseApi>({
      path: `/actions/${id}/enable`,
      method: 'PATCH',
      type: ContentType.Json,
      format: 'json',
      ...params,
    });
  /**
   * @description Disables one or multiple actions.
   *
   * @tags Actions
   * @name ActionsDisable
   * @summary Disable actions
   * @request PATCH:/actions
   */
  actionsDisable = (id: string, params: RequestParams = {}) =>
    this.http.request<PostResponseApi, GenericErrorResponseApi>({
      path: `/actions/${id}/disable`,
      method: 'PATCH',
      type: ContentType.Json,
      format: 'json',
      ...params,
    });
  /**
   * @description Patch an existing actions by ID.
   *
   * @tags Actions
   * @name ActionsPartialUpdate
   * @summary Update an actions
   * @request PATCH:/actions/
   */
  actionsPartialUpdate = (action: ActionRequestApi, params: RequestParams = {}) =>
    this.http.request<void, GenericErrorResponseApi>({
      path: `/actions/`,
      method: 'PATCH',
      body: action,
      type: ContentType.Json,
      ...params,
    });
  /**
   * @description Delete an action by ID.
   *
   * @tags Actions
   * @name ActionsDelete
   * @summary Delete an action
   * @request DELETE:/actions/{id}
   */
  actionsDelete = (id: string, params: RequestParams = {}) =>
    this.http.request<void, GenericErrorResponseApi>({
      path: `/actions/${id}`,
      method: 'DELETE',
      type: ContentType.Json,
      ...params,
    });
}
