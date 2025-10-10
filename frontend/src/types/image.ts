export interface ImageEditorParams {
  RedBrightness: number;
  GreenBrightness: number;
  BlueBrightness: number;
  Contrast: number;
  Negative: boolean;
  Order: string;
  VerticalMirror: boolean;
  HorizontalMirror: boolean;
  Magic: number;
  Filter: string;
  FilterSize: 3;
  Sigma: number;
  Interval: number;
  UnsharpMasking: number;
}

export interface ChannelImage {
  type: 'red' | 'green' | 'blue' | 'gray' | 'changes';
  url: string;
  histogram?: number[];
}

export interface BrightnessHistogram {
  Brightness: {
    Red: number[];
    Green: number[];
    Blue: number[];
    Gray: number[];
  };
}
