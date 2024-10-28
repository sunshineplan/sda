import { expect, test } from 'vitest'
import utils, { format, preprocess } from './utils'

test('Contains', () => {
    expect(utils.contains(['1', '2', '3'], '1')).toBeTruthy()
    expect(utils.contains(['1', '2', '3'], '4')).toBeFalsy()
})

test('Range', () => expect(utils.range('2', '5')).toStrictEqual(['2', '3', '4', '5']))

test('Sort', () => expect(utils.sort(['10', '2', '1', 'a'])).toStrictEqual(['1', '2', '10', 'a']))

test('Preprocess', () => {
    expect(preprocess('3\n1\n\n\n2\na\n\n\n')).toStrictEqual(['3', '1', '2', 'a'])
    expect(preprocess('3\n"1\n\n\n2"\n"a"\n\n\n')).toStrictEqual(['3', '1', '2', '"a"'])
})

test('Format', () => {
    expect(format(1, ['1'])).toBe('\nTotal 1 record\n\nresult:\n1')
    expect(format(2, ['1', '2'])).toBe('\nTotal 2 records\n\nresult:\n1\n2')
})
