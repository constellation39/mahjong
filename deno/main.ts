// 查表法得到向听数 https://zhuanlan.zhihu.com/p/31000381

// 对排序枚举
function EnumTiles(sum: number): Array<Number> {
  let ret: Array<Number> = [];
  for (const stack of QuantityEnum(sum)) {
    // console.log(stack);
    for (const tiles of DistanceEnum(stack)) {
      if (Valid(tiles)) {
        ret.push(Rebuild(tiles));
      }
    }
  }
  return ret;
}

// 对数量枚举
function QuantityEnum(sum: number): Array<number[]> {
  let ret: Array<number[]> = new Array();
  quantityEnum([], sum, ret);
  return ret;
}

function quantityEnum(stack: number[], sum: number, ret: Array<number[]>) {
  if (sum == 0) {
    ret.push(stack.reverse().concat());
    return;
  }

  for (let index = 1; index <= 4; index++) {
    if (sum >= index) {
      stack.push(index);
      quantityEnum(stack, sum - index, ret);
      stack.pop();
    }
  }
}

function DistanceEnum(stack: number[]): Array<number> {
  let ret: Array<number> = new Array();
  distanceEnum(
    Build(stack) | 0b1000 << ((stack.length - 1) * 4),
    0,
    (stack.length - 1) * 2,
    0,
    ret,
  );
  return ret;
}

function Build(stack: number[]): number {
  let ret = 0;
  let shift = 0;
  for (let i = 0; i < stack.length; i++, shift += 4) {
    ret |= (stack[i] - 1) << shift;
  }
  ret |= 0xF << shift;
  return ret;
}

function Valid(value: number): boolean {
  let level = 0;
  let continuous = 0, tempContinuous = 1;
  for (let shift = 0; (value >> shift) != 0xF; shift += 4) {
    let singleContinuous = (value >> (shift + 2)) & 3;
    if (singleContinuous < 2) {
      if (level >= 3) return false;
      tempContinuous += singleContinuous + 1;
    } else if (level < 3) {
      if (continuous + 2 + tempContinuous <= 9) {
        continuous += 2 + tempContinuous;
      } else {
        continuous = 0;
        tempContinuous = 1;
        level++;
      }
    }
  }
  return true;
}

// 对距离枚举
function distanceEnum(
  value: number,
  deep: number,
  deepCnt: number,
  index: number,
  ret: number[],
) {
  if (deep >= deepCnt) {
    ret.push(value);
    return;
  }
  for (let i = 1; i <= 3; i++) {
    if (i == 3) {
      distanceEnum(value | i << (deep * 4 + 2), deep + 1, deepCnt, 0, ret);
      continue;
    }
    distanceEnum(
      value | i << (deep * 4 + 2),
      deep + 1,
      deepCnt,
      index + i,
      ret,
    );
  }
}

interface RebuildItem {
  value: number;
  pos: number;
}

function Rebuild(value: number): number {
  let shift = 0;
  let cache: Array<RebuildItem> = [];

  while (value != 0xF) {
    if (((value >> shift) & 0b1000) == 0) {
      shift += 4;
      continue;
    }
    shift += 4;
    cache.push({ value: value & ((1 << shift) - 1), pos: shift });
    value >>= shift;
    shift = 0;
  }
  cache.sort((a: RebuildItem, b: RebuildItem) => a.value - b.value);
  for (let i = 0; i < cache.length; i++) {
    value = (value << cache[i].pos) | cache[i].value;
  }
  return value;
}

let cnt = 0;
for (let index = 2; index <= 14; index += 3) {
  for (const tiles of EnumTiles(index)) {
    cnt++;
    // console.log(tiles.toString(2));
    if (cnt % 10000 == 0) {
      console.log(cnt);
    }
  }
}
console.log(cnt);
