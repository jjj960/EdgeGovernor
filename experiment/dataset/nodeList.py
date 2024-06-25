
NodesList = {'node1': [64, 49512, 40, 1200, 220, [],[0.20, 0.10, 0.35, 0.40]],      #数据分别代表 [CPU总容量,内存总容量,带宽大小,磁盘总容量,磁盘IO速率,存在Pod列表,[已用CPU,已用内存,已用带宽,已用磁盘]]

                 'node2': [32, 32768, 30, 1000, 200, [],[0.30, 0.10, 0.42,0.20]],

                 'node3': [16, 24576, 30, 800, 180, [],[0.10, 0.10, 0.28,0.54]],

                 'node4': [64, 49512, 40, 1200, 220, [],[0.3, 0.20, 0.39,0.20]],

                 'node5': [32, 32768, 30, 1000, 200, [],[0.08, 0.07, 0.10,0.05]],

                 'node6': [16, 24576, 30, 800, 180, [],[0.10, 0.10, 0.30,0.45]],

                 'node7': [64, 49512, 40, 1200, 220, [],[0.30, 0.30, 0.16,0.06]],

                 'node8': [32, 32768, 30, 1000, 200, [],[0.20, 0.20, 0.3, 0.34]],

                 'node9': [16, 24576, 30, 800, 180, [],[0.10, 0.10, 0.33, 0.14]],

                 'node10': [64, 49512, 40, 1200, 220, [],[0.30, 0.20, 0.24, 0.09]],

                 'node11': [32, 32768, 30, 1000, 200, [],[0.20, 0.10, 0.35, 0.25]],

                 'node12': [16, 24576, 30, 800, 180, [],[0.10, 0.05, 0.19, 0.31]],

                 'node13': [64, 49512, 40, 1200, 220, [],[0.30, 0.03, 0.06, 0.04]],

                 'node14': [32, 32768, 30, 1000, 200, [],[0.10, 0.20, 0.29, 0.36]],

                 'node15': [16, 24576, 30, 800, 180, [],[0.10, 0.01, 0.15, 0.29]],

                 'node16': [16, 24576, 30, 800, 180, [],[0.20, 0.10, 0.15, 0.29]],

                 'node17': [64, 49512, 40, 1200, 220, [],[0.30, 0.10, 0.08, 0.09]],

                 'node18': [32, 32768, 30, 1000, 200, [],[0.20, 0.02, 0.02, 0.06]],

                 'node19': [16, 24576, 30, 800, 180, [],[0.10, 0.05, 0.22, 0.6]],

                 'node20': [64, 49512, 40, 1200, 220, [],[0.06, 0.10, 0.15, 0.45]],

                 'node21': [32, 32768, 30, 1000, 200, [],[0.30, 0.20, 0.24, 0.09]],

                 'node22': [16, 24576, 30, 800, 180, [],[0.09, 0.10, 0.4, 0.1]],

                 'node23': [64, 49512, 40, 1200, 220, [],[0.30, 0.20, 0.05, 0.1]],

                 'node24': [32, 32768, 30, 1000, 200, [],[0.05, 0.20, 0.05, 0.01]],

                 'node25': [16, 24576, 30, 800, 180, [],[0.10, 0.10, 0.14, 0.06]],

                 'node26': [64, 49512, 40, 1200, 220, [],[0.30, 0.20, 0.29, 0.08]],

                 'node27': [32, 32768, 30, 1000, 200, [],[0.20, 0.30, 0.4, 0.09]],

                 'node28': [16, 24576, 30, 800, 180, [],[0.08, 0.10, 0.09, 0.36]],

                 'node29': [64, 49512, 40, 1200, 220, [],[0.20, 0.30, 0.058, 0.25]],

                 'node30': [32, 32768, 30, 1000, 200, [],[0.20, 0.20, 0.4, 0.23]],

                 'node31': [16, 24576, 30, 800, 180, [],[0.10, 0.15, 0.05, 0.18]]}


# NodesList = {'node-15': [64, 48, 40, 1200, 220, [],[0.20, 0.10, 0.35, 0.40]],
#
#                  'node-14': [32, 32, 30, 1000, 200, [],[0.30, 0.10, 0.42,0.20]],
#
#                  'node-13': [16, 24, 30, 800, 180, [],[0.10, 0.10, 0.28,0.54]],
#
#                  'node-12': [64, 48, 40, 1200, 220, [],[0.3, 0.20, 0.39,0.20]],
#
#                  'node-11': [32, 32, 30, 1000, 200, [],[0.08, 0.07, 0.10,0.05]],
#
#                  'node-10': [16, 24, 30, 800, 180, [],[0.10, 0.10, 0.30,0.45]],
#
#                  'node-9': [64, 48, 40, 1200, 220, [],[0.30, 0.30, 0.16,0.06]],
#
#                  'node-8': [32, 32, 30, 1000, 200, [],[0.20, 0.20, 0.3, 0.34]],
#
#                  'node-7': [16, 24, 30, 800, 180, [],[0.10, 0.10, 0.33, 0.14]],
#
#                  'node-6': [64, 49, 40, 1200, 220, [],[0.30, 0.20, 0.24, 0.09]],
#
#                  'node-5': [32, 32, 30, 1000, 200, [],[0.20, 0.10, 0.35, 0.25]],
#
#                  'node-4': [16, 24, 30, 800, 180, [],[0.10, 0.05, 0.19, 0.31]],
#
#                  'node-3': [64, 48, 40, 1200, 220, [],[0.30, 0.03, 0.06, 0.04]],
#
#                  'node-2': [32, 32, 30, 1000, 200, [],[0.10, 0.20, 0.29, 0.36]],
#
#                  'node-1': [16, 24, 30, 800, 180, [],[0.10, 0.01, 0.15, 0.29]],
#
#                  'node0': [16, 24, 30, 800, 180, [],[0.20, 0.10, 0.15, 0.29]],
#
#                  'node1': [64, 48, 40, 1200, 220, [],[0.30, 0.10, 0.08, 0.09]],
#
#                  'node2': [32, 32, 30, 1000, 200, [],[0.20, 0.02, 0.02, 0.06]],
#
#                  'node3': [16, 24, 30, 800, 180, [],[0.10, 0.05, 0.22, 0.6]],
#
#                  'node4': [64, 48, 40, 1200, 220, [],[0.06, 0.10, 0.15, 0.45]],
#
#                  'node5': [32, 32, 30, 1000, 200, [],[0.30, 0.20, 0.24, 0.09]],
#
#                  'node6': [16, 24, 30, 800, 180, [],[0.09, 0.10, 0.4, 0.1]],
#
#                  'node7': [64, 48, 40, 1200, 220, [],[0.30, 0.20, 0.05, 0.1]],
#
#                  'node8': [32, 32, 30, 1000, 200, [],[0.05, 0.20, 0.05, 0.01]],
#
#                  'node9': [16, 24, 30, 800, 180, [],[0.10, 0.10, 0.14, 0.06]],
#
#                  'node10': [64, 48, 40, 1200, 220, [],[0.30, 0.20, 0.29, 0.08]],
#
#                  'node11': [32, 32, 30, 1000, 200, [],[0.20, 0.30, 0.4, 0.09]],
#
#                  'node12': [16, 24, 30, 800, 180, [],[0.08, 0.10, 0.09, 0.36]],
#
#                  'node13': [64, 48, 40, 1200, 220, [],[0.20, 0.30, 0.058, 0.25]],
#
#                  'node14': [32, 32, 30, 1000, 200, [],[0.20, 0.20, 0.4, 0.23]],
#
#                  'node15': [16, 24, 30, 800, 180, [],[0.10, 0.15, 0.05, 0.18]]}
#
