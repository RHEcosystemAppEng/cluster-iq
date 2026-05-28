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

import { ExpenseListResponseApi, ExpenseRequestApi, GenericErrorResponseApi, PostResponseApi } from './data-contracts';
import { ContentType, HttpClient, RequestParams } from './http-client';

export class Expenses<SecurityDataType = unknown> {
  http: HttpClient<SecurityDataType>;

  constructor(http: HttpClient<SecurityDataType>) {
    this.http = http;
  }

  /**
   * @description Paginated retrieval with optional filters.
   *
   * @tags Expenses
   * @name ExpensesList
   * @summary List expenses
   * @request GET:/expenses
   */
  expensesList = (
    query?: {
      /** Instance ID filter */
      instance_id?: string;
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
    this.http.request<ExpenseListResponseApi, GenericErrorResponseApi>({
      path: `/expenses`,
      method: 'GET',
      query: query,
      type: ContentType.Json,
      format: 'json',
      ...params,
    });
  /**
   * @description Create one or multiple expense records.
   *
   * @tags Expenses
   * @name ExpensesCreate
   * @summary Create expenses
   * @request POST:/expenses
   */
  expensesCreate = (expenses: ExpenseRequestApi[], params: RequestParams = {}) =>
    this.http.request<PostResponseApi, GenericErrorResponseApi>({
      path: `/expenses`,
      method: 'POST',
      body: expenses,
      type: ContentType.Json,
      format: 'json',
      ...params,
    });
}
