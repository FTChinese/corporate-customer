import { Injectable } from '@angular/core';
import { Subject } from 'rxjs';
import { AbstractControl } from '@angular/forms';

@Injectable({
  providedIn: 'root'
})
export class SearchService {

  private valueControlSource = new Subject<AbstractControl>();

  valueSubmitted$ = this.valueControlSource.asObservable();

  submitForm(control: AbstractControl) {
    this.valueControlSource.next(control);
  }
}
