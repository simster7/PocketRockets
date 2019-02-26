# -*- coding: utf-8 -*-
# snapshottest: v1 - https://goo.gl/zC4yUc
from __future__ import unicode_literals

from snapshottest import Snapshot


snapshots = Snapshot()

snapshots['TestGame.test_game_basic 1'] = "{'seats': [{'name': 'Simon', 'stack': 100, 'seat_number': 0}, {'name': 'Hersh', 'stack': 100, 'seat_number': 1}, {'name': 'Chien', 'stack': 100, 'seat_number': 2}, {'name': 'Jarry', 'stack': 100, 'seat_number': 3}, None, None, {'name': 'Grace', 'stack': 100, 'seat_number': 6}, {'name': 'Jason', 'stack': 100, 'seat_number': 7}, None], 'button_position': 0, 'small_blind': 1, 'big_blind': 2, 'game_state': None, 'test': True}"

snapshots['TestGame.test_game_basic 2'] = "{'players': [{'name': 'Simon', 'stack': 99, 'seat_number': 0}, {'name': 'Hersh', 'stack': 98, 'seat_number': 1}, {'name': 'Chien', 'stack': 100, 'seat_number': 2}, {'name': 'Jarry', 'stack': 100, 'seat_number': 3}, {'name': 'Grace', 'stack': 100, 'seat_number': 6}, {'name': 'Jason', 'stack': 100, 'seat_number': 7}], 'bet_vector': [1, 2, 0, 0, 0, 0], 'fold_vector': [False, False, False, False, False, False], 'pot': 0, 'acting_player': 2, 'leading_player': 1, 'deck': [{'card_id': 0}, {'card_id': 1}, {'card_id': 2}, {'card_id': 3}, {'card_id': 4}, {'card_id': 5}, {'card_id': 6}, {'card_id': 7}, {'card_id': 8}, {'card_id': 9}, {'card_id': 10}, {'card_id': 11}, {'card_id': 12}, {'card_id': 13}, {'card_id': 14}, {'card_id': 15}, {'card_id': 16}, {'card_id': 17}, {'card_id': 18}, {'card_id': 19}, {'card_id': 20}, {'card_id': 21}, {'card_id': 22}, {'card_id': 23}, {'card_id': 24}, {'card_id': 25}, {'card_id': 26}, {'card_id': 27}, {'card_id': 28}, {'card_id': 29}, {'card_id': 30}, {'card_id': 31}, {'card_id': 32}, {'card_id': 33}, {'card_id': 34}, {'card_id': 35}, {'card_id': 36}, {'card_id': 37}, {'card_id': 38}, {'card_id': 39}, {'card_id': 40}, {'card_id': 41}, {'card_id': 42}, {'card_id': 43}, {'card_id': 44}, {'card_id': 45}, {'card_id': 46}, {'card_id': 47}, {'card_id': 48}, {'card_id': 49}, {'card_id': 50}, {'card_id': 51}], 'round': 0}"

snapshots['TestGame.test_game_basic 3'] = "{'bet_round': 0, 'lead_action': {'action': <Actions.bet: 3>, 'value': 7}, 'lead_player': {'name': 'Chien', 'stack': 93, 'seat_number': 2}, 'current_players': [{'name': 'Simon', 'stack': 99, 'seat_number': 0}, {'name': 'Hersh', 'stack': 98, 'seat_number': 1}, {'name': 'Chien', 'stack': 93, 'seat_number': 2}, {'name': 'Jarry', 'stack': 93, 'seat_number': 3}, {'name': 'Grace', 'stack': 100, 'seat_number': 6}, {'name': 'Jason', 'stack': 100, 'seat_number': 7}], 'player_cards': None, 'community_cards': [], 'acting_player': {'name': 'Jason', 'stack': 100, 'seat_number': 7}, 'end_game': None}"

snapshots['TestGame.test_game_basic 4'] = "{'players': [{'name': 'Simon', 'stack': 93, 'seat_number': 0}, {'name': 'Hersh', 'stack': 98, 'seat_number': 1}, {'name': 'Chien', 'stack': 93, 'seat_number': 2}, {'name': 'Jarry', 'stack': 93, 'seat_number': 3}, {'name': 'Grace', 'stack': 100, 'seat_number': 6}, {'name': 'Jason', 'stack': 93, 'seat_number': 7}], 'bet_vector': [0, 0, 0, 0, 0, 0], 'fold_vector': [False, True, False, False, True, False], 'pot': 30, 'acting_player': 0, 'leading_player': 0, 'deck': [{'card_id': 0}, {'card_id': 1}, {'card_id': 2}, {'card_id': 3}, {'card_id': 4}, {'card_id': 5}, {'card_id': 6}, {'card_id': 7}, {'card_id': 8}, {'card_id': 9}, {'card_id': 10}, {'card_id': 11}, {'card_id': 12}, {'card_id': 13}, {'card_id': 14}, {'card_id': 15}, {'card_id': 16}, {'card_id': 17}, {'card_id': 18}, {'card_id': 19}, {'card_id': 20}, {'card_id': 21}, {'card_id': 22}, {'card_id': 23}, {'card_id': 24}, {'card_id': 25}, {'card_id': 26}, {'card_id': 27}, {'card_id': 28}, {'card_id': 29}, {'card_id': 30}, {'card_id': 31}, {'card_id': 32}, {'card_id': 33}, {'card_id': 34}, {'card_id': 35}, {'card_id': 36}, {'card_id': 37}, {'card_id': 38}, {'card_id': 39}, {'card_id': 40}, {'card_id': 41}, {'card_id': 42}, {'card_id': 43}, {'card_id': 44}, {'card_id': 45}, {'card_id': 46}, {'card_id': 47}, {'card_id': 48}, {'card_id': 49}, {'card_id': 50}, {'card_id': 51}], 'round': 1}"

snapshots['TestGame.test_game_basic 5'] = "{'bet_round': 1, 'lead_action': {'action': <Actions.check: 0>, 'value': None}, 'lead_player': {'name': 'Simon', 'stack': 93, 'seat_number': 0}, 'current_players': [{'name': 'Simon', 'stack': 93, 'seat_number': 0}, {'name': 'Hersh', 'stack': 98, 'seat_number': 1}, {'name': 'Chien', 'stack': 93, 'seat_number': 2}, {'name': 'Jarry', 'stack': 93, 'seat_number': 3}, {'name': 'Grace', 'stack': 100, 'seat_number': 6}, {'name': 'Jason', 'stack': 93, 'seat_number': 7}], 'player_cards': [{'card_id': 0}, {'card_id': 6}], 'community_cards': [{'card_id': 13}, {'card_id': 14}, {'card_id': 15}], 'acting_player': {'name': 'Simon', 'stack': 93, 'seat_number': 0}, 'end_game': None}"

snapshots['TestGame.test_game_basic 6'] = "{'bet_round': 1, 'lead_action': {'action': <Actions.bet: 3>, 'value': 10}, 'lead_player': {'name': 'Jarry', 'stack': 83, 'seat_number': 3}, 'current_players': [{'name': 'Simon', 'stack': 93, 'seat_number': 0}, {'name': 'Hersh', 'stack': 98, 'seat_number': 1}, {'name': 'Chien', 'stack': 93, 'seat_number': 2}, {'name': 'Jarry', 'stack': 83, 'seat_number': 3}, {'name': 'Grace', 'stack': 100, 'seat_number': 6}, {'name': 'Jason', 'stack': 93, 'seat_number': 7}], 'player_cards': None, 'community_cards': [{'card_id': 13}, {'card_id': 14}, {'card_id': 15}], 'acting_player': {'name': 'Simon', 'stack': 93, 'seat_number': 0}, 'end_game': None}"

snapshots['TestGame.test_game_basic 7'] = "{'players': [{'name': 'Simon', 'stack': 83, 'seat_number': 0}, {'name': 'Hersh', 'stack': 98, 'seat_number': 1}, {'name': 'Chien', 'stack': 83, 'seat_number': 2}, {'name': 'Jarry', 'stack': 83, 'seat_number': 3}, {'name': 'Grace', 'stack': 100, 'seat_number': 6}, {'name': 'Jason', 'stack': 93, 'seat_number': 7}], 'bet_vector': [0, 0, 0, 0, 0, 0], 'fold_vector': [False, True, False, False, True, True], 'pot': 60, 'acting_player': 0, 'leading_player': 0, 'deck': [{'card_id': 0}, {'card_id': 1}, {'card_id': 2}, {'card_id': 3}, {'card_id': 4}, {'card_id': 5}, {'card_id': 6}, {'card_id': 7}, {'card_id': 8}, {'card_id': 9}, {'card_id': 10}, {'card_id': 11}, {'card_id': 12}, {'card_id': 13}, {'card_id': 14}, {'card_id': 15}, {'card_id': 16}, {'card_id': 17}, {'card_id': 18}, {'card_id': 19}, {'card_id': 20}, {'card_id': 21}, {'card_id': 22}, {'card_id': 23}, {'card_id': 24}, {'card_id': 25}, {'card_id': 26}, {'card_id': 27}, {'card_id': 28}, {'card_id': 29}, {'card_id': 30}, {'card_id': 31}, {'card_id': 32}, {'card_id': 33}, {'card_id': 34}, {'card_id': 35}, {'card_id': 36}, {'card_id': 37}, {'card_id': 38}, {'card_id': 39}, {'card_id': 40}, {'card_id': 41}, {'card_id': 42}, {'card_id': 43}, {'card_id': 44}, {'card_id': 45}, {'card_id': 46}, {'card_id': 47}, {'card_id': 48}, {'card_id': 49}, {'card_id': 50}, {'card_id': 51}], 'round': 2}"

snapshots['TestGame.test_game_basic 8'] = "{'players': [{'name': 'Simon', 'stack': 83, 'seat_number': 0}, {'name': 'Hersh', 'stack': 98, 'seat_number': 1}, {'name': 'Chien', 'stack': 83, 'seat_number': 2}, {'name': 'Jarry', 'stack': 83, 'seat_number': 3}, {'name': 'Grace', 'stack': 100, 'seat_number': 6}, {'name': 'Jason', 'stack': 93, 'seat_number': 7}], 'bet_vector': [0, 0, 0, 0, 0, 0], 'fold_vector': [False, True, False, False, True, True], 'pot': 60, 'acting_player': 0, 'leading_player': 0, 'deck': [{'card_id': 0}, {'card_id': 1}, {'card_id': 2}, {'card_id': 3}, {'card_id': 4}, {'card_id': 5}, {'card_id': 6}, {'card_id': 7}, {'card_id': 8}, {'card_id': 9}, {'card_id': 10}, {'card_id': 11}, {'card_id': 12}, {'card_id': 13}, {'card_id': 14}, {'card_id': 15}, {'card_id': 16}, {'card_id': 17}, {'card_id': 18}, {'card_id': 19}, {'card_id': 20}, {'card_id': 21}, {'card_id': 22}, {'card_id': 23}, {'card_id': 24}, {'card_id': 25}, {'card_id': 26}, {'card_id': 27}, {'card_id': 28}, {'card_id': 29}, {'card_id': 30}, {'card_id': 31}, {'card_id': 32}, {'card_id': 33}, {'card_id': 34}, {'card_id': 35}, {'card_id': 36}, {'card_id': 37}, {'card_id': 38}, {'card_id': 39}, {'card_id': 40}, {'card_id': 41}, {'card_id': 42}, {'card_id': 43}, {'card_id': 44}, {'card_id': 45}, {'card_id': 46}, {'card_id': 47}, {'card_id': 48}, {'card_id': 49}, {'card_id': 50}, {'card_id': 51}], 'round': 3}"

snapshots['TestGame.test_game_basic 9'] = "{'bet_round': 3, 'lead_action': {'action': <Actions.check: 0>, 'value': None}, 'lead_player': {'name': 'Simon', 'stack': 83, 'seat_number': 0}, 'current_players': [{'name': 'Simon', 'stack': 83, 'seat_number': 0}, {'name': 'Hersh', 'stack': 98, 'seat_number': 1}, {'name': 'Chien', 'stack': 83, 'seat_number': 2}, {'name': 'Jarry', 'stack': 83, 'seat_number': 3}, {'name': 'Grace', 'stack': 100, 'seat_number': 6}, {'name': 'Jason', 'stack': 93, 'seat_number': 7}], 'player_cards': [{'card_id': 3}, {'card_id': 9}], 'community_cards': [{'card_id': 13}, {'card_id': 14}, {'card_id': 15}, {'card_id': 17}, {'card_id': 19}], 'acting_player': {'name': 'Simon', 'stack': 83, 'seat_number': 0}, 'end_game': None}"

snapshots['TestGame.test_game_basic 10'] = "{'players': [{'name': 'Simon', 'stack': 73, 'seat_number': 0}, {'name': 'Hersh', 'stack': 98, 'seat_number': 1}, {'name': 'Chien', 'stack': 83, 'seat_number': 2}, {'name': 'Jarry', 'stack': 73, 'seat_number': 3}, {'name': 'Grace', 'stack': 100, 'seat_number': 6}, {'name': 'Jason', 'stack': 93, 'seat_number': 7}], 'bet_vector': [0, 0, 0, 0, 0, 0], 'fold_vector': [False, True, True, False, True, True], 'pot': 80, 'acting_player': 0, 'leading_player': 0, 'deck': [{'card_id': 0}, {'card_id': 1}, {'card_id': 2}, {'card_id': 3}, {'card_id': 4}, {'card_id': 5}, {'card_id': 6}, {'card_id': 7}, {'card_id': 8}, {'card_id': 9}, {'card_id': 10}, {'card_id': 11}, {'card_id': 12}, {'card_id': 13}, {'card_id': 14}, {'card_id': 15}, {'card_id': 16}, {'card_id': 17}, {'card_id': 18}, {'card_id': 19}, {'card_id': 20}, {'card_id': 21}, {'card_id': 22}, {'card_id': 23}, {'card_id': 24}, {'card_id': 25}, {'card_id': 26}, {'card_id': 27}, {'card_id': 28}, {'card_id': 29}, {'card_id': 30}, {'card_id': 31}, {'card_id': 32}, {'card_id': 33}, {'card_id': 34}, {'card_id': 35}, {'card_id': 36}, {'card_id': 37}, {'card_id': 38}, {'card_id': 39}, {'card_id': 40}, {'card_id': 41}, {'card_id': 42}, {'card_id': 43}, {'card_id': 44}, {'card_id': 45}, {'card_id': 46}, {'card_id': 47}, {'card_id': 48}, {'card_id': 49}, {'card_id': 50}, {'card_id': 51}], 'round': 4}"

snapshots['TestGame.test_game_basic 11'] = "{'players': [{'name': 'Simon', 'stack': 72, 'seat_number': 0}, {'name': 'Hersh', 'stack': 96, 'seat_number': 1}, {'name': 'Chien', 'stack': 83, 'seat_number': 2}, {'name': 'Jarry', 'stack': 73, 'seat_number': 3}, {'name': 'Grace', 'stack': 100, 'seat_number': 6}, {'name': 'Jason', 'stack': 93, 'seat_number': 7}], 'bet_vector': [1, 2, 0, 0, 0, 0], 'fold_vector': [False, False, False, False, False, False], 'pot': 0, 'acting_player': 2, 'leading_player': 1, 'deck': [{'card_id': 0}, {'card_id': 1}, {'card_id': 2}, {'card_id': 3}, {'card_id': 4}, {'card_id': 5}, {'card_id': 6}, {'card_id': 7}, {'card_id': 8}, {'card_id': 9}, {'card_id': 10}, {'card_id': 11}, {'card_id': 12}, {'card_id': 13}, {'card_id': 14}, {'card_id': 15}, {'card_id': 16}, {'card_id': 17}, {'card_id': 18}, {'card_id': 19}, {'card_id': 20}, {'card_id': 21}, {'card_id': 22}, {'card_id': 23}, {'card_id': 24}, {'card_id': 25}, {'card_id': 26}, {'card_id': 27}, {'card_id': 28}, {'card_id': 29}, {'card_id': 30}, {'card_id': 31}, {'card_id': 32}, {'card_id': 33}, {'card_id': 34}, {'card_id': 35}, {'card_id': 36}, {'card_id': 37}, {'card_id': 38}, {'card_id': 39}, {'card_id': 40}, {'card_id': 41}, {'card_id': 42}, {'card_id': 43}, {'card_id': 44}, {'card_id': 45}, {'card_id': 46}, {'card_id': 47}, {'card_id': 48}, {'card_id': 49}, {'card_id': 50}, {'card_id': 51}], 'round': 0}"
