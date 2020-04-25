import { Injectable } from '@angular/core';
import { Subject } from 'rxjs';
import { RequestError } from '../models/request-result';

@Injectable()
export class FormService {

  private formSubmittedSource = new Subject<string>();
  private errorsSource = new Subject<RequestError>();

  formSubmitted$ = this.formSubmittedSource.asObservable();
  errorReceived$ = this.errorsSource.asObservable();

  // Pass form data to parent host.
  submit(data: string) {
    this.formSubmittedSource.next(data);
  }

  // Parent host send error messages to form.
  sendError(err: RequestError) {
    this.errorsSource.next(err);
  }
}
