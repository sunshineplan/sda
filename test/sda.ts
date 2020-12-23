import * as sda from '../src/sda'

test('Check Duplicates', () => {
  expect(new sda.chkDuplicates(['1', '2', '3']).run()).toStrictEqual({})
  expect(new sda.chkDuplicates(['2', '1', '2', '3', '1', '2']).run()).toStrictEqual({ 2: 3, 1: 2 })
})

test('Remove Duplicates', () => expect(new sda.rmDuplicates(['1', '2', '3', '1']).run()).toStrictEqual(['1', '2', '3']))

test('Compare', () => {
  const data1 = ['1', '2', '3']
  const data2 = ['1', '2', '3', '4']
  expect(new sda.compare(data1, data2).run()).toStrictEqual(['1', '2', '3'])
  expect(new sda.compare(data1, data2).run()).toStrictEqual(['1', '2', '3'])
  expect(new sda.compare(data1, data2, 'diff').run()).toHaveLength(0)
  expect(new sda.compare(data2, data1, 'diff').run()).toStrictEqual(['4'])
})

test('Check Consecutive', () => {
  expect(new sda.chkConsecutive(["5", "1", "3"]).run()).toStrictEqual(['2', '4'])
  expect(new sda.chkConsecutive(["3", "2", "1"]).run()).toHaveLength(0)
  expect(new sda.chkConsecutive(["1", "2", "a"]).run()).toStrictEqual(['!Error!'])
})
