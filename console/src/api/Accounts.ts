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
  AccountListResponseApi,
  AccountRequestApi,
  AccountResponseApi,
  ClusterListResponseApi,
  GenericErrorResponseApi,
  InstanceListResponseApi,
  PostResponseApi,
} from './data-contracts';
import { ContentType, HttpClient, RequestParams } from './http-client';

export class Accounts<SecurityDataType = unknown> {
  http: HttpClient<SecurityDataType>;

  constructor(http: HttpClient<SecurityDataType>) {
    this.http = http;
  }

  /**
   * @description Paginated retrieval of accounts with optional filters.
   *
   * @tags Accounts
   * @name AccountsList
   * @summary List accounts
   * @request GET:/accounts
   */
  accountsList = (
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
      /** Cloud provider */
      provider?: string;
    },
    params: RequestParams = {}
  ) =>
    this.http.request<AccountListResponseApi, GenericErrorResponseApi>({
      path: `/accounts`,
      method: 'GET',
      query: query,
      type: ContentType.Json,
      format: 'json',
      ...params,
    });
  /**
   * @description Create one or multiple accounts from the request body.
   *
   * @tags Accounts
   * @name AccountsCreate
   * @summary Create accounts
   * @request POST:/accounts
   */
  accountsCreate = (accounts: AccountRequestApi[], params: RequestParams = {}) =>
    this.http.request<PostResponseApi, GenericErrorResponseApi>({
      path: `/accounts`,
      method: 'POST',
      body: accounts,
      type: ContentType.Json,
      format: 'json',
      ...params,
    });
  /**
   * @description Return a single account resource.
   *
   * @tags Accounts
   * @name AccountsDetail
   * @summary Get an account by ID
   * @request GET:/accounts/{id}
   */
  accountsDetail = (id: string, params: RequestParams = {}) =>
    this.http.request<AccountResponseApi, GenericErrorResponseApi>({
      path: `/accounts/${id}`,
      method: 'GET',
      type: ContentType.Json,
      format: 'json',
      ...params,
    });
  /**
   * @description Delete an account resource by its ID.
   *
   * @tags Accounts
   * @name AccountsDelete
   * @summary Delete an account
   * @request DELETE:/accounts/{id}
   */
  accountsDelete = (id: string, params: RequestParams = {}) =>
    this.http.request<void, GenericErrorResponseApi>({
      path: `/accounts/${id}`,
      method: 'DELETE',
      type: ContentType.Json,
      ...params,
    });
  /**
   * @description Patch an existing account by ID.
   *
   * @tags Accounts
   * @name AccountsPartialUpdate
   * @summary Update an account
   * @request PATCH:/accounts/{id}
   */
  accountsPartialUpdate = (id: string, account: AccountRequestApi, params: RequestParams = {}) =>
    this.http.request<void, void>({
      path: `/accounts/${id}`,
      method: 'PATCH',
      body: account,
      type: ContentType.Json,
      ...params,
    });
  /**
   * @description Return the clusters associated with the specified account.
   *
   * @tags Accounts
   * @name ClustersList
   * @summary List clusters by account ID
   * @request GET:/accounts/{id}/clusters
   */
  clustersList = (id: string, params: RequestParams = {}) =>
    this.http.request<ClusterListResponseApi, GenericErrorResponseApi>({
      path: `/accounts/${id}/clusters`,
      method: 'GET',
      type: ContentType.Json,
      format: 'json',
      ...params,
    });
  /**
   * @description Return the list of instances for expense update for a specified account.
   *
   * @tags Accounts
   * @name ExpenseUpdateInstancesList
   * @summary List expense update instances
   * @request GET:/accounts/{id}/expense_update_instances
   */
  expenseUpdateInstancesList = (id: string, params: RequestParams = {}) =>
    this.http.request<InstanceListResponseApi, GenericErrorResponseApi>({
      path: `/accounts/${id}/expense_update_instances`,
      method: 'GET',
      type: ContentType.Json,
      format: 'json',
      ...params,
    });
}
