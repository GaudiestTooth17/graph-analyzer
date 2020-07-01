import sys


def main():
    with open(sys.argv[1], 'r') as in_file:
        with open(sys.argv[1][:-4] + '.csv', 'w') as out_file:
            for line in in_file:
                out_file.write(','.join(line.split(' ')))


if __name__ == '__main__':
    main()
