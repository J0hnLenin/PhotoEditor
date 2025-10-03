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
  Sigma: Number;
}

export interface ChannelImage {
  type: 'red' | 'green' | 'blue' | 'gray';
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

// Значения по умолчанию для сброса
export const DEFAULT_EDITOR_PARAMS: ImageEditorParams = {
  RedBrightness: 100,
  GreenBrightness: 100,
  BlueBrightness: 100,
  Contrast: 100,
  Negative: false,
  Order: 'RGB',
  VerticalMirror: false,
  HorizontalMirror: false,
  Magic: 0,
  Filter:  'none',
  FilterSize: 3,
  Sigma: 0,
};