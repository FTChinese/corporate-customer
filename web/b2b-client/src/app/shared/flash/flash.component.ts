import { Component, Input, Output, EventEmitter} from '@angular/core';
import { AlertType } from '../alert';

@Component({
  selector: 'app-flash',
  templateUrl: './flash.component.html'
})
export class FlashComponent  {
  @Input() message: string;
  @Input() type: AlertType = 'danger';
  @Input() dismissible = true;

  @Output() dismiss = new EventEmitter<void>();

  close() {
    this.dismiss.emit();
  }
}
